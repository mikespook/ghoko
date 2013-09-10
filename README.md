GHoKo
=====

Every GitLab project can trigger a web server whenever the repo is pushed to. 
Web Hooks can be used to update an external issue tracker, trigger CI builds,
update a backup mirror, or even deploy to your production server.

GitLab will send POST request with commits information on every push.

GHoKo is a web application that listens to post-hooks from GitLab, scripted
by Lua and written by Golang.

Installing
----------

	go get github.com/mikespook/ghoko

Usage
-----

	ghoko -h
