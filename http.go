// Copyright 2013 Xing Xing <mikespook@gmail.com>.
// All rights reserved.
// Use of this source code is governed by a commercial
// license that can be found in the LICENSE file.

package ghoko

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/mikespook/golib/idgen"
	"github.com/mikespook/golib/iptpool"
	"github.com/mikespook/golib/log"
)

type Handler struct {
	scriptPath string
	secret     string
	idgen      idgen.IdGen
	iptPool    *iptpool.IptPool
}

func New(scriptPath, secret string) (h *Handler) {
	h = &Handler{
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

func writeAndLog(w http.ResponseWriter, r *http.Request, status int, data []byte) {
	log.Messagef("%s %s \"%s\" %d", r.RemoteAddr, r.Method, r.URL.String(), status)
	w.WriteHeader(status)
	if data != nil {
		if _, err := w.Write(data); err != nil {
			log.Errorf("%s %s \"%s\" %d \"%s\"", r.RemoteAddr, r.Method, r.URL.String(), status, err)
		}
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u, err := url.ParseRequestURI(r.RequestURI)
	if err != nil {
		writeAndLog(w, r, http.StatusInternalServerError, []byte(err.Error()))
		return
	}
	if u.Query().Get("_secret") != h.secret {
		writeAndLog(w, r, http.StatusForbidden, nil)
		return
	}
	hoko, err := newHook(h, w, r)
	if err != nil {
		writeAndLog(w, r, http.StatusInternalServerError, []byte(err.Error()))
		return
	}
	status, data := hoko.exec()
	writeAndLog(w, r, status, data)
}

func (h *Handler) post(uri string, params Params) ([]byte, error) {
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

func (h *Handler) postJson(uri string, params Params) ([]byte, error) {
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

func (h *Handler) get(uri string) ([]byte, error) {
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

func (h *Handler) call(id, name string, params Params) error {
	ipt := h.iptPool.Get()
	defer h.iptPool.Put(ipt)
	ipt.Bind("Id", id)
	return ipt.Exec(name, params)
}
