// Copyright 2013 Xing Xing <mikespook@gmail.com>.
// All rights reserved.
// Use of this source code is governed by a commercial
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"github.com/mikespook/ghoko"
	"github.com/mikespook/golib/log"
	"github.com/mikespook/golib/pid"
	"github.com/mikespook/golib/signal"
	"net/http"
	"os"
	"path"
	"time"
)

var (
	addr       string
	scriptPath string
	secret     string
	tlsCert    string
	tlsKey     string
	pidfile    string
)

func init() {
	if !flag.Parsed() {
		flag.StringVar(&addr, "addr", ":3080", "Address of HTTP service")
		flag.StringVar(&scriptPath, "script", "./", "Path of lua files")
		flag.StringVar(&secret, "secret", "", "Secret token")
		flag.StringVar(&tlsCert, "tls-cert", "", "TLS cert file")
		flag.StringVar(&tlsKey, "tls-key", "", "TLS key file")
		flag.StringVar(&pidfile, "pid", "", "PID file")
		flag.Parse()
	}
	log.InitWithFlag()
}

func main() {
	log.Messagef("Starting: addr=%q script=%q", addr, scriptPath)
	if pidfile != "" {
		if p, err := pid.New(pidfile); err != nil {
			log.Error(err)
		} else {
			defer func() {
				if err := p.Close(); err != nil {
					log.Error(err)
				}
			}()
			log.Messagef("PID: %d file=%q", p.Pid, pidfile)
		}
	}
	defer func() {
		log.Message("Exited!")
		time.Sleep(time.Millisecond * 100)
	}()

	// Begin
	p := path.Clean(scriptPath)
	ghk := ghoko.New(p, secret)
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
