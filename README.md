GHoKo
=====

GHoKo is a web application that listens to post-hooks from [GitLab][gitlab]
and [GitHub][github], scripted by Lua and written in [Golang][golang].

For GitLab
----------
Every [GitLab][gitlab] project can trigger a web server whenever the repo is
pushed to. Web Hooks can be used to update an external issue tracker, 
trigger [CI][ci] builds, update a backup mirror, or even deploy to your 
production server.

GitLab will send POST request with commits information on every push.

For GitHub
----------
Every [GitHub][github] repository has the option to communicate with a web
server whenever the repository is pushed to. These "WebHooks" can be used to 
update an external issue tracker, trigger [CI][ci] builds, update a backup
mirror, or even deploy to your production server.

When a push is made to your repository, we'll POST to your URL with a payload
of JSON-encoded data about the push and the commits it contained. 

Dependency
==========

 * [mikespook/golib][4]
 * [aarzilli/golua][5]
 * [stevedonovan/luar][6]

Installing
==========

	go get github.com/mikespook/ghoko

Usage
=====

Service
-------

Executing:

	ghoko -h

For help information:

	Usage of ./ghoko:
		-addr=":8080": Address of http service
		-log="": log to write (empty for STDOUT)
		-log-level="all": log level ('error', 'warning', 'message', 'debug',
		'all' and 'none' are combined with '|')
		-main="gitlab": Main hosted repository
		-script="./": Path of lua files
		-secret="": Secret token

Scripting
---------

GHoKo use Lua as the scripting language. GHoKo will pass the Request into
Lua script as a user data. The structure of Request is [here][gitlab-req]
for gitlab (sign in needed) and [here][github-req] for github.

You can use this user data in the [Lua script][demo].

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
[4]: https://github.com/mikespook/golib
[5]: https://github.com/aarzilli/golua
[6]: https://github.com/stevedonovan/luar
[demo]: https://github.com/mikespook/ghoko/blob/master/foobar.lua
[blog]: http://mikespook.com
[twitter]: http://twitter.com/mikespook
[github-req]: https://help.github.com/articles/post-receive-hooks
[gitlab-req]: http://demo.gitlab.com/help/web_hooks
