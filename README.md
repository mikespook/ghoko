GHoKo
=====

[![Build Status][travis-img]][travis]

GHoKo is a web application that listens to web-hooks, scripted by Lua and
written in [Golang][golang]. Web-hooks usually is used for [CI][ci].

Dependency
==========

 * [mikespook/golib][golib]
 * [aarzilli/golua][golua]
 * [stevedonovan/luar][luar]

Installing
==========

	go get github.com/mikespook/ghoko

Usage
=====

Service
-------

Executing following command:

	$ ${fullpath}/ghoko -h

We will get:

	Usage of ./ghoko:
		-addr=":8080": Address of http service
		-defualt="gitlab": Default code hosting site
		-log="": log to write (empty for STDOUT)
		-log-level="all": log level ('error', 'warning', 'message', 'debug', 'all' and 'none' are combined with '|')
		-pid="": PID file
		-script="./": Path of lua files
		-secret="": Secret token
		-tls-cert="": TLS cert file
		-tls-key="": TLS key file
		

The pattern of hook URL is 

	${schema}://${addr}/${default}/${hook}?secret=${secret}&${params}

$schema could be HTTP or HTTPS either. When both two `tls-*` flags were
specified correctly, The HTTPS will be used.

$params can be used for passing custom values into script.

Scripting
---------

GHoKo use Lua as the scripting language. GHoKo will pass the Request into
Lua script as a user data.

You can use this user data in the [Lua script][demo].

Web Hook
--------

To set GitLab's web hook: Your repo --> settrings --> Web Hooks.

To set GitHub's web hook is a little more complicated.
Following: Your repo --> Settings --> Service Hooks --> WebHook URLs.

Here is an example for gitlab ([gitlab.lua][gitlab-lua]):

	http://192.168.1.100/gitlab?secret=phrase

or for github ([github.lua][github-lua]):

	http://192.168.1.100/gitlab?secret=phrase

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
[travis-img]: https://travis-ci.org/mikespook/z-node.png?branch=master
[travis]: https://travis-ci.org/mikespook/z-node
