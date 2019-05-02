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

SupportedOperatingSystems = darwin linux windows

.PHONY: all depEnsure test clean cleanAll build install buildSupported $(SupportedOperatingSystems)

.DEFAULT_GOAL := build

build: depEnsure
	CGO_ENABLED=0 go build -o build/${APP}-${GOHOSTOS}-${GOHOSTARCH} -ldflags "-s -w" -a -installsuffix cgo .

${GOPATH}/bin/dep:
	go get -v github.com/golang/dep/cmd/dep

depEnsure: ${GOPATH}/bin/dep
	${GOPATH}/bin/dep ensure -v -vendor-only

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

${SupportedOperatingSystems}: depEnsure
	CGO_ENABLED=0 GOOS=$@ GOARCH=amd64 go build -o build/${APP}-$@-amd64${EXT} -ldflags "-s -w" -a -installsuffix cgo .

all: $(SupportedOperatingSystems)

logs/unittests.log:
	-mkdir -p logs
	-go test -v ./functions -tags unit 2>&1 | tee logs/unittests.log

logs/functionaltests.log:
	-mkdir -p logs
	-go test -test.timeout=0 -v ./functionaltests 2>&1 | tee logs/functionaltests.log

${GOPATH}/bin/go-junit-report:
	go get github.com/jstemmer/go-junit-report

logs/junit-report.xml: ${GOPATH}/bin/go-junit-report logs/unittests.log logs/functionaltests.log
	cat logs/unittests.log logs/functionaltests.log | ${GOPATH}/bin/go-junit-report | tee logs/junit-report.xml

test : logs/junit-report.xml
