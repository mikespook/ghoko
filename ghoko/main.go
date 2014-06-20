// Copyright 2013 Xing Xing <mikespook@gmail.com>.
// All rights reserved.
// Use of this source code is governed by a commercial
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/mikespook/ghoko"
	"github.com/mikespook/golib/log"
	"github.com/mikespook/golib/pid"
	"github.com/mikespook/golib/signal"
)

var (
	addr       string
	scriptPath string
	secret     string
	tlsCert    string
	tlsKey     string
	pidFile    string
	rootUrl    string
)

func init() {
	if !flag.Parsed() {
		flag.StringVar(&addr, "addr", ":3080", "Address of HTTP service")
		flag.StringVar(&scriptPath, "script", path.Dir(os.Args[0]), "Path of lua files")
		flag.StringVar(&secret, "secret", "", "Secret token")
		flag.StringVar(&tlsCert, "tls-cert", "", "TLS cert file")
		flag.StringVar(&tlsKey, "tls-key", "", "TLS key file")
		flag.StringVar(&pidFile, "pid", "", "PID file")
		flag.StringVar(&rootUrl, "root", "/", "Root path of URL")
		flag.Parse()
	}
	log.InitWithFlag()
}

func main() {
	log.Messagef("Starting: webhook=%q script=%q", ghoko.CallbackUrl(tlsCert, tlsKey, addr, rootUrl), scriptPath)
	if pidFile != "" {
		if p, err := pid.New(pidFile); err != nil {
			log.Error(err)
		} else {
			defer func() {
				if err := p.Close(); err != nil {
					log.Error(err)
				}
			}()
			log.Messagef("PID: %d file=%q", p.Pid, pidFile)
		}
	}
	defer func() {
		log.Message("Exited!")
		time.Sleep(time.Millisecond * 100)
	}()

	// Begin
	p := path.Clean(scriptPath)
	ghk := ghoko.New(p, secret, rootUrl)
	go func() {
		defer func() {
			if err := signal.Send(os.Getpid(), os.Interrupt); err != nil {
				panic(err)
			}
		}()
		var err error
		if tlsCert == "" || tlsKey == "" {
			err = http.ListenAndServe(addr, ghk)
		} else {
			err = http.ListenAndServeTLS(addr, tlsCert, tlsKey, ghk)
		}
		if err != nil {
			log.Error(err)
		}
	}()
	// End

	sh := signal.NewHandler()
	sh.Bind(os.Interrupt, func() bool { return true })
	sh.Loop()
}
