// Copyright 2013 Xing Xing <mikespook@gmail.com>.
// All rights reserved.
// Use of this source code is governed by a commercial
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"github.com/mikespook/golib/log"
	"github.com/mikespook/golib/signal"
	"net"
	"os"
	"syscall"
	"time"
)

const (
	GITLAB = "gitlab"
	GITHUB = "github"
)

var (
	addr       = flag.String("addr", ":8080", "Address of http service")
	scriptPath = flag.String("script", "./", "Path of lua files")
	secret     = flag.String("secret", "", "Secret token")
	mainHosting = flag.String("main", "gitlab", "Main hosted repository")
)

func init() {
	if !flag.Parsed() {
		flag.Parse()
	}
	log.Flag()
	if *mainHosting != GITHUB {
		*mainHosting = GITLAB
	}
}

func main() {
	log.Message("Starting...")

	defer func() {
		log.Message("Exit.")
		time.Sleep(time.Second)
	}()
	p := *scriptPath
	if p[len(p)-1] == 47 {
		p = p[:len(p)-1]
	}
	hook := NewHook(*addr, p, *secret, *mainHosting)
	go func() {
		if e := hook.Serve(); e != nil {
			if _, ok := e.(*net.OpError); !ok {
				log.Error(e)
			}
		}
	}()
	defer hook.Close()
	sh := signal.NewHandler()
	sh.Bind(os.Interrupt, func() bool { return true })
	sh.Bind(syscall.SIGUSR1, func() bool {
		if e := hook.Close(); e != nil {
			log.Error(e)
		}
		var attr os.ProcAttr
		attr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
		attr.Sys = &syscall.SysProcAttr{}
		_, e := os.StartProcess(os.Args[0], os.Args, &attr)
		if e != nil {
			log.Error(e)
		}
		return true
	})
	sh.Loop()
}
