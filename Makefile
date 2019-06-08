GOOS ?= linux
ARCH ?= amd64

VERSION := $(shell git rev-parse --short=8 HEAD)

all: build test

build: _deps
	VERSION=$(VERSION) hack/build.sh

test:
	go test ./...

run:
	go run cmd/main.go \
		-kubeconfig=$$KUBECONFIG \
		-logtostderr=true \
		-v=4

_deps:
	dep version || curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	dep ensure -v
	dep prune -v
	find vendor -name '*_test.go' -delete
