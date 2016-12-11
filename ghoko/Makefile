VERSION=1.0.0
PLATFORM=$(shell uname -s -m)
BUILT=$(shell date)
build:
	go build -ldflags='-X "main.version=${VERSION}" -X "main.platform=$(PLATFORM)" -X "main.built=$(BUILT)"'
clean:
	go clean
