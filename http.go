// Copyright 2013 Xing Xing <mikespook@gmail.com>.
// All rights reserved.
// Use of this source code is governed by a commercial
// license that can be found in the LICENSE file.

package ghoko

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mikespook/golib/idgen"
	"github.com/mikespook/golib/iptpool"
	"github.com/mikespook/golib/log"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

var (
	ErrAccessDeny       = errors.New("Access Deny")
	ErrMethodNotAllowed = errors.New("Method Not Allowed")
	ErrSyncNeeded       = errors.New("`sync` param needed")
)

type ghokoHandler struct {
	scriptPath string
	secret     string
	idgen      idgen.IdGen
	iptPool    *iptpool.IptPool
}

func New(scriptPath, secret string) (h *ghokoHandler) {
	h = &ghokoHandler{
		scriptPath: scriptPath,
		secret:     secret,
		idgen:      idgen.NewObjectId(),
		iptPool:    iptpool.NewIptPool(NewLuaIpt),
	}
	h.iptPool.OnCreate = func(ipt iptpool.ScriptIpt) error {
		ipt.Init(h.scriptPath)
		ipt.Bind("Call", h.call)
		ipt.Bind("Secret", h.secret)
		return nil
	}
	return h
}

func (h *ghokoHandler) Close() error {
	errstr := ""
	emap := h.iptPool.Free()
	for k, err := range emap {
		errstr = fmt.Sprintf("%s[%s]: %s\n", errstr, k, err)
	}
	if errstr != "" {
		return errors.New(errstr)
	}
	return nil
}

func (h *ghokoHandler) verify(p url.Values) bool {
	if h.secret == "" {
		return true
	}
	return h.secret == p.Get("secret")
}

func (h *ghokoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fallthrough
	case "POST":
	default:
		log.Errorf("[%s] %s \"%s: %s\"", r.RemoteAddr, r.RequestURI, ErrMethodNotAllowed, r.Method)
		http.Error(w, ErrMethodNotAllowed.Error(), 405)
		return
	}
	u, err := url.Parse(r.RequestURI)
	if err != nil {
		log.Errorf("[%s] %s \"%s\"", r.RemoteAddr, r.RequestURI, err)
		http.Error(w, err.Error(), 500)
		return
	}
	p := u.Query()
	if !h.verify(p) { // verify secret token
		log.Errorf("[%s] %s \"%s\"", r.RemoteAddr, r.RequestURI, ErrAccessDeny)
		http.Error(w, ErrAccessDeny.Error(), 403)
		return
	}
	p.Del("secret")
	params := make(Params)
	params.AddValues(p)
	if r.Method == "POST" {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Errorf("[%s] %s \"%s\"", r.RemoteAddr, r.RequestURI, err)
			http.Error(w, err.Error(), 500)
			return
		}
		defer r.Body.Close()
		if err := params.AddJSON(data); err != nil {
			log.Errorf("[%s] %s \"%s\"", r.RemoteAddr, r.RequestURI, err)
			http.Error(w, err.Error(), 500)
			return
		}
	}
	name := path.Base(u.Path)
	id := h.idgen.Id().(string)
	f := func(sync bool) {
		ipt := h.iptPool.Get()
		defer h.iptPool.Put(ipt)
		ipt.Bind("Id", id)
		ipt.Bind("WriteBody", func(str string) (err error) {
			if !sync {
				return ErrSyncNeeded
			}
			_, err = w.Write([]byte(str))
			return
		})
		ipt.Bind("WriteHeader", func(status int) error {
			if !sync {
				return ErrSyncNeeded
			}
			w.WriteHeader(status)
			return nil
		})

		if err := ipt.Exec(name, params); err != nil {
			log.Errorf("[%s] %s \"%s\"", r.RemoteAddr,
				r.RequestURI, err.Error())
			if sync {
				http.Error(w, err.Error(), 500)
			}
			return
		}
		log.Messagef("[%s] %s \"Success\"", r.RemoteAddr,
			r.RequestURI)
	}

	if p.Get("sync") == "true" {
		f(true)
		w.Header().Set("Ghoko-Id", id)
	} else {
		go f(false)
		if _, err := w.Write([]byte(fmt.Sprintf("\"%s\"", id))); err != nil {
			log.Errorf("[%s] %s %s \"%s\"", r.RemoteAddr,
				r.RequestURI, id, err)
		}
	}
}

func (h *ghokoHandler) call(id, name string, params Params) (err error) {
	ipt := h.iptPool.Get()
	defer h.iptPool.Put(ipt)
	ipt.Bind("Id", id)
	return ipt.Exec(name, params)
}
