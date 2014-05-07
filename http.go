// Copyright 2013 Xing Xing <mikespook@gmail.com>.
// All rights reserved.
// Use of this source code is governed by a commercial
// license that can be found in the LICENSE file.

package ghoko

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/mikespook/golib/idgen"
	"github.com/mikespook/golib/iptpool"
	"github.com/mikespook/golib/log"
)

var (
	ErrForbidden        = errors.New("Access Deny")
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
		ipt.Bind("Get", h.get)
		ipt.Bind("PostJSON", h.postJson)
		ipt.Bind("Post", h.post)
		ipt.Bind("Secret", h.secret)
		return nil
	}
	return h
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
		h.write(w, r, http.StatusMethodNotAllowed, ErrMethodNotAllowed.Error())
		return
	}
	u, err := url.Parse(r.RequestURI)
	if err != nil {
		log.Errorf("[%s] %s \"%s\"", r.RemoteAddr, r.RequestURI, err)
		h.write(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	p := u.Query()
	if !h.verify(p) { // verify secret token
		log.Errorf("[%s] %s \"%s\"", r.RemoteAddr, r.RequestURI, ErrForbidden)
		h.write(w, r, http.StatusForbidden, ErrForbidden.Error())
		return
	}
	p.Del("secret")
	params := make(Params)
	params.AddValues(p)
	if r.Method == "POST" {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Errorf("[%s] %s \"%s\"", r.RemoteAddr, r.RequestURI, err)
			h.write(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		defer r.Body.Close()
		if err := params.AddJSON(data); err != nil {
			log.Errorf("[%s] %s \"%s\"", r.RemoteAddr, r.RequestURI, err)
			h.write(w, r, http.StatusInternalServerError, err.Error())
			return
		}
	}
	name := path.Base(u.Path)
	var id string
	if params["_id"] == nil {
		id = h.idgen.Id().(string)
	} else {
		id = params["_id"].(string)
	}
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
		h.write(w, r, http.StatusOK, id)
	}
}

func (h *ghokoHandler) post(uri string, params Params) ([]byte, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Add("secret", h.secret)
	resp, err := http.PostForm(u.String(), params.Values())
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(string(body))
	}
	return body, nil
}

func (h *ghokoHandler) postJson(uri string, params Params) ([]byte, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Add("secret", h.secret)
	j, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(u.String(), "application/json", bytes.NewBuffer(j))
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(string(body))
	}
	return body, nil
}

func (h *ghokoHandler) get(uri string) ([]byte, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Add("secret", h.secret)
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(string(body))
	}
	return body, nil
}

func (h *ghokoHandler) call(id, name string, params Params) error {
	ipt := h.iptPool.Get()
	defer h.iptPool.Put(ipt)
	ipt.Bind("Id", id)
	return ipt.Exec(name, params)
}

func (h *ghokoHandler) write(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	accept := r.Header.Get("Accept")
	if strings.Contains(accept, "application/json") {
		if err := h.writeJson(w, status, data); err != nil {
			log.Errorf("[%s] %s \"%s\"", r.RemoteAddr, r.RequestURI, err)
		}
		return
	}
	w.WriteHeader(status)
	if _, err := w.Write([]byte(fmt.Sprintf("%s", data))); err != nil {
		log.Errorf("[%s] %s \"%s\"", r.RemoteAddr, r.RequestURI, err)
	}
}

func (h *ghokoHandler) writeJson(w http.ResponseWriter, status int, data interface{}) error {
	content, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.WriteHeader(status)
	_, err = w.Write(content)
	return err
}
