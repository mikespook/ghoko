package ghoko

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type hook struct {
	id      string
	isJson  bool
	isSync  bool
	w       http.ResponseWriter
	r       *http.Request
	params  Params
	name    string
	handler *Handler
}

func newHook(handler *Handler, w http.ResponseWriter, r *http.Request) (*hook, error) {
	id := r.Header.Get("Ghoko-Id")
	if id == "" {
		id = handler.idgen.Id().(string)
	}
	if !strings.HasPrefix(r.URL.Path, handler.rootUrl) {
		return nil, ErrNotFound
	}
	name := strings.TrimPrefix(r.URL.Path, handler.rootUrl)
	h := &hook{
		w:       w,
		r:       r,
		params:  make(Params),
		isJson:  strings.Contains(r.Header.Get("Content-Type"), "json"),
		isSync:  r.Header.Get("Ghoko-Sync") == "true",
		name:    name,
		handler: handler,
		id:      id,
	}
	if h.isJson {
		u, err := url.ParseRequestURI(r.RequestURI)
		if err != nil {
			return nil, err
		}
		h.params.AddValues(u.Query())
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		defer r.Body.Close()
		if err := h.params.AddJSON(data); err != nil {
			return nil, err
		}
	} else {
		if err := r.ParseForm(); err != nil {
			return nil, err
		}
		h.params.AddValues(r.Form)
	}
	return h, nil
}

func (h *hook) exec() (int, []byte) {
	f := func() (int, []byte, error) {
		ipt := h.handler.iptPool.Get()
		defer h.handler.iptPool.Put(ipt)
		var buf bytes.Buffer
		var status int
		ipt.Bind("Id", h.id)
		ipt.Bind("WriteBody", func(str string) error {
			if !h.isSync {
				return ErrSyncNeeded
			}
			_, err := buf.WriteString(str)
			return err
		})
		ipt.Bind("WriteHeader", func(s int) error {
			if !h.isSync {
				return ErrSyncNeeded
			}
			status = s
			return nil
		})

		if err := ipt.Exec(h.name, h.params); err != nil {
			if !h.isSync {
				writeAndLogError(nil, h.r, err)
			}
			return http.StatusInternalServerError, nil, err
		}
		return http.StatusOK, buf.Bytes(), nil
	}

	if h.isSync {
		h.w.Header().Set("Ghoko-Id", h.id)
		status, data, err := f()
		if err != nil {
			return http.StatusInternalServerError, []byte(err.Error())
		}
		return status, data
	}
	go f()
	return http.StatusOK, h.data(h.id)
}

func (h *hook) data(data string) []byte {
	if h.isJson {
		buf := bytes.NewBufferString("\"")
		buf.WriteString(data)
		buf.WriteString("\"")
		return buf.Bytes()
	}
	return []byte(data)
}
