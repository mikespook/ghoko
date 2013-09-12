// Copyright 2013 Xing Xing <mikespook@gmail.com>.
// All rights reserved.
// Use of this source code is governed by a commercial
// license that can be found in the LICENSE file.

package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mikespook/golib/iptpool"
	"github.com/mikespook/golib/log"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"path"
)

var (
	ErrAccessDeny = errors.New("Access Deny")
)

type httpServer struct {
	conn       net.Listener
	srv        *http.Server
	iptPool    *iptpool.IptPool
	secret     string
	scriptPath string
	hosting    string
}

func NewHook(addr, scriptPath, secret, hosting string) (srv *httpServer) {
	srv = &httpServer{
		srv:        &http.Server{Addr: addr},
		iptPool:    iptpool.NewIptPool(NewLuaIpt),
		scriptPath: scriptPath,
		secret:     secret,
		hosting:    hosting,
	}
	return
}

func (s *httpServer) SetTLS(certFile, keyFile string) (err error) {
	s.srv.TLSConfig = &tls.Config{}
	s.srv.TLSConfig.NextProtos = []string{"http/1.1"}
	s.srv.TLSConfig.Certificates = make([]tls.Certificate, 1)
	s.srv.TLSConfig.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	return
}

func (s *httpServer) Serve() (err error) {
	s.conn, err = net.Listen("tcp", s.srv.Addr)
	if err != nil {
		return
	}
	if s.srv.TLSConfig != nil {
		s.conn = tls.NewListener(s.conn, s.srv.TLSConfig)
	}
	s.iptPool.OnCreate = func(ipt iptpool.ScriptIpt) error {
		ipt.Init(s.scriptPath)
		return nil
	}
	http.HandleFunc("/", s.handler)
	return s.srv.Serve(s.conn)
}

func (s *httpServer) Close() error {
	errstr := ""
	emap := s.iptPool.Free()
	if n := len(emap); n > 0 {
		for k, err := range emap {
			errstr = fmt.Sprintf("%s[%s]: %s\n", errstr, k, err)
		}
	}
	s.conn.Close()
	if errstr != "" {
		return errors.New(errstr)
	}
	return nil
}

func (s *httpServer) handler(w http.ResponseWriter, r *http.Request) {
	host, name, params, req, err := s.prepare(r)
	if err != nil {
		log.Errorf("[%s] %s \"%s\"", r.RemoteAddr, r.RequestURI, err)
		http.Error(w, err.Error(), err.Errno())
		return
	}
	go func() {
		ipt := s.iptPool.Get()
		defer s.iptPool.Put(ipt)
		if req != nil {
			ipt.Bind("Host", host)
			switch r := req.(type) {
			case GitHubRequest:
				ipt.Bind("Request", r)
			case GitLabRequest:
				ipt.Bind("Request", r)
			}
		}
		if err := ipt.Exec(name, params); err != nil {
			log.Errorf("[%s] %s \"%s\"", r.RemoteAddr,
				r.RequestURI, err.Error())
			return
		}
		log.Messagef("[%s] %s \"Success\"", r.RemoteAddr,
			r.RequestURI)
	}()
	w.WriteHeader(200)
}

func (s *httpServer) prepare(r *http.Request) (host, name string, params url.Values, obj interface{}, gerr *ghokoErr) {
	u, err := url.Parse(r.RequestURI)
	if err != nil {
		gerr = NewError(err.Error(), 500)
		return
	}
	params = u.Query()
	if s.secret != params.Get("secret") { // verify secret token
		gerr = NewError(ErrAccessDeny.Error(), 403)
		return
	}
	params.Del("secret")
	name = path.Base(u.Path)
	host = params.Get("host")
	if host == "" {
		host = s.hosting
	}
	params.Del("default")
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			gerr = NewError(err.Error(), 500)
			return
		}
		defer r.Body.Close()
		switch host {
		case GITLAB:
			var req GitLabRequest
			if err = json.Unmarshal(body, &req); err != nil {
				gerr = NewError(err.Error(), 500)
				return
			}
			obj = req
		case GITHUB:
			var req GitHubRequest
			if err = json.Unmarshal(body, &req); err != nil {
				gerr = NewError(err.Error(), 500)
				return
			}
			obj = req
		}
	}
	return
}
