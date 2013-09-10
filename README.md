GHoKo
=====

Every [GitLab][1] project can trigger a web server whenever the repo is pushed to. 
Web Hooks can be used to update an external issue tracker, trigger [CI][2] builds,
update a backup mirror, or even deploy to your production server.

GitLab will send POST request with commits information on every push.

GHoKo is a web application that listens to post-hooks from GitLab, scripted
by Lua and written in [Golang][3].

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
		-addr=":8080": address of http service
		-log="": log to write (empty for STDOUT)
		-log-level="all": log level ('error', 'warning', 'message', 'debug', 'all' and 'none' are combined with '|')
		-script="./": script path
		-secret="": secret token

Scripting
---------

GHoKo use Lua as the scripting language. GHoKo will pass the gitlab's [Request][7] into Lua script as a user data.
The Request's struct (fake) is here:

	Request: {
		Before:		string,
  		After:		string,
  		Ref:		string,
  		UserId:		int,
		UserName:	string,
  		Repository {
			Name: 		string,
			Url:		string,
			Description:	string,
			Homepage:	string,
  		},
		Commits: [
    		Commit {
				Id:			string,
				Message:	string,
				Timestamp:	time.Time,
				Url:		string,
				Author: {
					Name:	string,
					Email:	string,
				},
    		},
  		],
		TotalCommitsCount: int,
	}

You can use this user data in the [Lua script][8].

Authors
=======

 * Xing Xing <mikespook@gmail.com> [Blog][blog] [@Twitter][twitter]

Open Source
===========

See LICENSE for more information.

[1]: http://www.gitlab.com
[2]: http://en.wikipedia.org/wiki/Continuous_integration
[3]: http://golang.org
[4]: https://github.com/mikespook/golib
[5]: https://github.com/aarzilli/golua
[6]: https://github.com/stevedonovan/luar
[7]: https://github.com/mikespook/ghoko/blob/master/types.go
[8]: https://github.com/mikespook/ghoko/blob/master/foobar.lua
[blog]: http://mikespook.com
[twitter]: http://twitter.com/mikespook
