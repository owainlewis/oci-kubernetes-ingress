GOOS ?= linux
ARCH ?= amd64

VERSION := $(shell git rev-parse --short=8 HEAD)

all: deps build test

build:
	VERSION=$(VERSION) hack/build.sh

test:
	go test ./...

run:
	go run cmd/main.go \
		-kubeconfig=$$KUBECONFIG \
		-v=4

deps:
	dep version || curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	dep ensure -v
	dep prune -v
	find vendor -name '*_test.go' -delete
