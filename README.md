GHoKo
=====

[![Build Status][travis-img]][travis]

GHoKo is a web application that listens to web-hooks, scripted by Lua and
written in [Golang][golang]. Web-hooks usually are used for [CI][ci], 
[Automated testing][auto-testing] or something else.

For example, GHoKo has been used in many browser based game projects for
CI and system operations.

Dependency
==========

 * [mikespook/golib][golib]
 * [aarzilli/golua][golua]
 * [stevedonovan/luar][luar]
 * [liblua5.1-0-dev][liblua] for Ubuntu

Installing
==========

All useful scripts were put at the directory [shell][shell].

Befor building, the proper lua librarie must be installed.
E.g. Ubuntu 14.04, it is `liblua5.1-0-dev`.

Then:

	go get github.com/mikespook/ghoko/ghoko

The ghoko library implement http.Handler and can be embeded
into other projects.

	go get github.com/mikespook/ghoko

Have a fun!

Usage
=====

Service
-------

Executing following command:

	$ ${fullpath}/ghoko -h

Some help information:

	Usage of ./ghoko:
		-addr=":8080": Address of http service
		-defualt="gitlab": Default code hosting site
		-log="": log to write (empty for STDOUT)
		-log-level="all": log level ('error', 'warning', 'message', 'debug', 
			'all' and 'none' are combined with '|')
		-pid="": PID file
		-root="/": Root path of URL
		-script="./": Path of lua files
		-secret="": Secret token
		-tls-cert="": TLS cert file
		-tls-key="": TLS key file
		

The pattern of hook URL is 

	${schema}://${addr}/${root}/${hook}?_secret=${secret}&${params}

`$schema` could be HTTP or HTTPS either. When both two `tls-*` flags were
specified correctly, The HTTPS will be used.

You can set root path of URL through `root` flag.

Eg. `script` was set to `/ghoko`. And if `root` is `/hook/v1`, the request
`http://127.0.0.1:3080/hook/v1/foo/bar` will evaluate `/ghoko/foo/bar.lua`.
If `root` was set to `/hook` and requesting the same URL, 
`/ghoko/v1/foo/bar.lua` will be evaluated.

`$params` can be used for passing custom values into script through URL. 
HTTP method, POST is also accepted. If `Content-Type` in the request header
contains `json`, it means passing enconded JSON data through POST-Body.
Otherwise, it is a common post with form data.

All of them will combine into a global variable `ghoko.Params`, it can
be used in Lua scripts.

Usually, GHoKo evaluates lua scripts asynchronous. `Ghoko-Sync` is a magic 
header for requesting ghoko in synchronized way. When it is equal 
`ture`(string), two functions `ghoko.WriteBody` and `ghoko.WriteHeader`
can be used for response data and HTTP status to HTTP clients.

Another magic header is `GHoKo-Id`. It tells ghoko do not generate ID
but using client specified one.

Scripting
---------

GHoKo use Lua as the scripting language. GHoKo will pass the Request into
Lua script as a user data.

You can use user data in the [Lua script][demo].

Following variables and functions can be called in Lua:

 * ghoko.Id - Every request has a global unique Id
 * ghoko.Params - Params passed by URL\POST-BODY(JSON format)
 * ghoko.Call(id, name, params) - Call lua script and pass params to it
 * ghoko.Debug(msg)/ghoko.Debugf(format, msg) - Output debug infomations
 * ghoko.Message(msg)/ghoko.Messagef(format, msg) - Output message infomations
 * ghoko.Warning(msg)/ghoko.Warningf(format, msg) - Output warning infomations
 * ghoko.Error(err)/ghoko.Errorf(format, msg) - Output error infomations
 * ghoko.Write(msg) - Write something to HTTP clients (sync only)
 * ghoko.WriteHeader(status) - Assign HTTP status (sync only)
 * ghoko.Get(url) - GET a remote url, `_secret` will be passed
 * ghoko.PostJSON(url, params) - POST to a remote url with JSON encoded params
 * ghoko.Post(url, params) - POST to a remote url with a form

Web Hook
--------

To set GitLab's web hook: Your repo --> settrings --> Web Hooks.

To set GitHub's web hook is a little more complicated.
Following: Your repo --> Settings --> Service Hooks --> WebHook URLs.

Here is an example for gitlab ([gitlab.lua][gitlab-lua]):

	http://192.168.1.100/gitlab?_secret=phrase

or for github ([github.lua][github-lua]):

	http://192.168.1.100/github?_secret=phrase

We have writen demo scripts for you. The scripts will print the repo and commits's informations.

Authors
=======

 * Xing Xing <mikespook@gmail.com> [Blog][blog] [@Twitter][twitter]

Open Source
===========

See LICENSE for more information.

[gitlab]: http://www.gitlab.com
[github]: http://www.github.com
[ci]: http://en.wikipedia.org/wiki/Continuous_integration
[golang]: http://golang.org
[golib]: https://github.com/mikespook/golib
[golua]: https://github.com/aarzilli/golua
[luar]: https://github.com/stevedonovan/luar
[demo]: https://github.com/mikespook/ghoko/blob/master/foobar.lua
[blog]: http://mikespook.com
[twitter]: http://twitter.com/mikespook
[github-req]: https://help.github.com/articles/post-receive-hooks
[gitlab-req]: http://demo.gitlab.com/help/web_hooks
[rhodecode]: https://rhodecode.com/
[bitbucket]: https://bitbucket.org/
[github-lua]: https://github.com/mikespook/ghoko/blob/master/github.lua
[gitlab-lua]: https://github.com/mikespook/ghoko/blob/master/gitlab.lua
[travis-img]: https://travis-ci.org/mikespook/ghoko.png?branch=master
[travis]: https://travis-ci.org/mikespook/ghoko
[auto-testing]: http://en.wikipedia.org/wiki/Test_automation
[shell]: https://github.com/mikespook/ghoko/tree/master/shell  
[liblua]: http://packages.ubuntu.com/trusty/liblua5.1-0-dev
