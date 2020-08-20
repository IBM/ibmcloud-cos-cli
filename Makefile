APP:=$(notdir $(CURDIR))
GOPATH:=$(shell go env GOPATH)
GOHOSTOS:=$(shell go env GOHOSTOS)
GOHOSTARCH:=$(shell go env GOHOSTARCH)

BUILD_DATE:="$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")"
GIT_COMMIT_SHA:="$(shell git rev-parse --abbrev-ref HEAD -->/dev/null)"
GIT_REMOTE_URL:="$(shell git config --get remote.origin.url 2>/dev/null)"
BUILD_DATE:="$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")"

# Dependencies package
PACKAGES:=vendor

# Jenkins vars. Set to `unknown` if the variable is not yet defined
BUILD_ID?=unknown
BUILD_NUMBER?=unknown
SHELL:=/bin/bash

SupportedOperatingSystems = darwin linux windows ppc64le

.PHONY: all test clean cleanAll build install buildSupported $(SupportedOperatingSystems) set_env_vars buildWithMod vendorWithMod

.DEFAULT_GOAL := build

build:
	CGO_ENABLED=0 go build -o build/${APP}-${GOHOSTOS}-${GOHOSTARCH} -ldflags "-s -w" -a -installsuffix cgo .

clean:
	-go clean
	-go clean -testcache
	-rm -rf build
	-rm -rf $(APP)-*
	-rm -rf logs

cleanAll : clean
	-rm -rf vendor

install:
	ibmcloud plugin install build/${APP}-${GOHOSTOS}-${GOHOSTARCH}

osx : darwin
windows : EXT := .exe

buildSupported: $(SupportedOperatingSystems)

${SupportedOperatingSystems}:
	if [ "$@" != "ppc64le" ]; then\
	    CGO_ENABLED=0 GOOS=$@ GOARCH=amd64 go build -o build/${APP}-$@-amd64${EXT} -ldflags "-s -w" -a -installsuffix cgo . ;\
	else\
	    CGO_ENABLED=0 GOOS=linux GOARCH=$@ go build -o build/${APP}-linux-$@ -ldflags "-s -w" -a -installsuffix cgo . ;\
	fi

all: $(SupportedOperatingSystems)

test :
	go test -v -tags unit ./functions

#### experimental cli go mod support ###
## populate vendor folder
vendorWithMod :
	go mod vendor

## Build using go mod
buildWithMod: vendorWithMod
	GOFLAGS="-mod=vendor $$GOFLAGS" \
	CGO_ENABLED=0 \
	go build -o build/${APP}-${GOHOSTOS}-${GOHOSTARCH} -ldflags "-s -w" -a -installsuffix cgo .
