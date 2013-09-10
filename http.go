// Copyright 2013 Xing Xing <mikespook@gmail.com>.
// All rights reserved.
// Use of this source code is governed by a commercial
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"github.com/mikespook/golib/iptpool"
	"github.com/mikespook/golib/log"
	"net"
	"net/http"
	"net/url"
	"io/ioutil"
		"encoding/json"
)

var (
	ErrAccessDeny = errors.New("Access Deny")
	ErrPostOnly = errors.New("POST method only")
)

type httpServer struct {
	conn       net.Listener
	srv        *http.Server
	iptPool    *iptpool.IptPool
	secret     string
	scriptPath string
}

func NewHook(addr, scriptPath, secret string) (srv *httpServer) {
	srv = &httpServer{
		srv:        &http.Server{Addr: addr},
		iptPool:    iptpool.NewIptPool(NewLuaIpt),
		scriptPath: scriptPath,
		secret:     secret,
	}
	return
}

func (s *httpServer) Serve() (err error) {
	s.conn, err = net.Listen("tcp", s.srv.Addr)
	if err != nil {
		return
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
	if r.Method != "POST" { // only post method permited
		log.Errorf("[%s] %s \"%s\"", r.RemoteAddr, r.RequestURI,
			ErrPostOnly)
		http.Error(w, ErrPostOnly.Error(), 500)
	}

	p, err := url.Parse(r.RequestURI)
	if err != nil {
		log.Errorf("[%s] %s \"%s\"", r.RemoteAddr, r.RequestURI, err)
		http.Error(w, err.Error(), 500)
		return
	}

	if s.secret != p.Query().Get("secret") { // verify secret token
		log.Errorf("[%s] %s \"%s\"", r.RemoteAddr, r.RequestURI, ErrAccessDeny)
		http.Error(w, ErrAccessDeny.Error(), 403)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("[%s] %s \"%s\"", r.RemoteAddr, r.RequestURI, err)
		http.Error(w, err.Error(), 500)
		return
	}
	defer r.Body.Close()

	go func() {
		ipt := s.iptPool.Get()
		defer s.iptPool.Put(ipt)
		var postReq PostRequest
		err := json.Unmarshal(body, &postReq)
		if err != nil {
			log.Errorf("[%s] %s \"%s\"", r.RemoteAddr,
				r.RequestURI, err.Error())
			return
		}
		ipt.Bind("Request", &postReq)
		if err := ipt.Exec(p.Path, nil); err != nil {
			log.Errorf("[%s] %s \"%s\"", r.RemoteAddr,
				r.RequestURI, err.Error())
			return
		}
		log.Messagef("[%s] %s \"Success\"", r.RemoteAddr,
			r.RequestURI)
	}()
	w.WriteHeader(200)
}
