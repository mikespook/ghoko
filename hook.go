package ghoko

import (
	"net/http"
	"path"
	"strings"

	"github.com/mikespook/golib/idgen"
)

type hoko struct {
	id     string
	isJson bool
	sync   bool
	w      http.ResponseWriter
	params Params
	name   string
	secret string
}

func newHoko(id string, w http.ResponseWriter, r *http.Request) (*hoko, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	id := r.Header.Get("Ghoko-Id")
	if id == "" {
		id = h.idgen.Id().(string)
	}
	h := &hoko{
		w:      w,
		params: make(Params),
		isJson: strings.Contains(r.Header.Get("Content-Type"), "json"),
		sync:   r.Header.Get("Sync") == "true",
		name:   path.Base(r.URL.Path),
	}
	h.params.AddValues(r.Form)
	if h.isJson {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		defer r.Body.Close()
		if err := params.AddJSON(data); err != nil {
			h.write(w, r, http.StatusInternalServerError, err.Error())
			return nil, err
		}
	}
	return h, nil
}

func (h *hoko) forbidden(secret string) bool {
	return secret != "" && secret != h.secret
}

func (h *hoko) exec() (int, []byte) {
	f := func(sync bool) {
		ipt := h.iptPool.Get()
		defer h.iptPool.Put(ipt)
		ipt.Bind("Id", h.id)
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
				h.write(w, r, http.StatusInternalServerError, err.Error())
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

func (h *hoko) write() error {
	if h.isJson {
		return h.writeJson()
	}
	return h.writeText()
}

func (h *hoko) writeJson() error {
	content, err := json.Marshal(h.data)
	if err != nil {
		return err
	}
	w.WriteHeader(h.status)
	_, err = w.Write(content)
	return err
}

func (h *hoko) writeText() error {
	w.WriteHeader(h.status)
	_, err = w.Write(h.content)
	return err
}
