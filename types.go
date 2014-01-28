// Copyright 2013 Xing Xing <mikespook@gmail.com>.
// All rights reserved.
// Use of this source code is governed by a commercial
// license that can be found in the LICENSE file.

package ghoko

import (
	"encoding/json"
	"github.com/stevedonovan/luar"
	"net/url"
)

type Params luar.Map

func (p Params) AddValues(values url.Values) {
	for k, v := range values {
		p[k] = v
	}
}

func (p Params) AddJSON(data []byte) (err error) {
	var tmp luar.Map
	if err = json.Unmarshal(data, &tmp); err != nil {
		return
	}
	for k, v := range tmp {
		p[k] = v
	}
	return
}
