GOOS?=linux
GOARCH?=amd64
GOBIN:=$(shell pwd)/.bin

GIT_REPO=$(shell git config --get remote.origin.url)
GIT_COMMIT=$(shell git rev-parse --short HEAD)

VERSION := $(shell git rev-parse --short=8 HEAD)

VERSION_LD_FLAGS=
COMPILE_OUTPUT?=controller

all: build

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -v -ldflags="-s -w $(VERSION_LD_FLAGS)" -a -installsuffix cgo  -o ${COMPILE_OUTPUT} ./cmd

.PHONY: run
run:
	go run cmd/main.go
