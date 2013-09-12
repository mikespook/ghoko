GHoKo
=====

GHoKo is a web application that listens to web-hooks, scripted by Lua and
written in [Golang][golang]. Web-hooks usually is used for [CI][ci].

Currently GHoKo officially supports the following code hosting sites:

 * [GitLab][gitlab]
 * [GitHub][github]

And the following are being planned:
 
 * [BitBucket][bitbucket]
 * [RHodeCode][rhodecode]

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
		-default="gitlab": Default code hosting site
		-log="": log to write (empty for STDOUT)
		-log-level="all": log level ('error', 'warning', 'message', 'debug',
		'all' and 'none' are combined with '|')
		-main="gitlab": Main hosted repository
		-script="./": Path of lua files
		-secret="": Secret token
		-tls-cert="": TLS cert file
		-tls-key="": TLS key file

The pattern of hook URL is 

	${schema}://${addr}/${default}/${hook}?secret=${secret}&default=${default}&${params}

$schema could be HTTP or HTTPS either. When both two `tls-*` flags were
specified correctly, The HTTPS will be used.

The flag `default` will be used as default code hosting site. If it is not
specified, "gitlab" as a default. You can also pass a url query parameter
`default` to override this value. When the value matches "gitlab", 
[the hook request of gitlab][gitlab-req] will be passed into the scripts;
if it matches "github", [github's request][github-req] will be passed;
otherwise, `ghoko.Request` variable will not be set in the scripts.

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

Here is an example for gitlab:

	http://192.168.1.100/ci?secret=phrase&default=gitlab&test=true

or for github:

	http://192.168.1.100/test?secret=phrase&default=gitlab&callback=someurl

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
