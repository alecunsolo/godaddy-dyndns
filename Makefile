# see https://gist.github.com/serinth/16391e360692f6a000e5a10382d1148c
SERVICE  ?= $(shell basename `go list`)
VERSION  ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || cat $(PWD)/.version 2> /dev/null || echo v0)
PACKAGE  ?= $(shell go list)
PACKAGES ?= $(shell go list ./...)
FILES    ?= $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: default
default: help

.PHONY: help
help:   ## show this help
	@echo 'usage: make [target] ...'
	@echo ''
	@echo 'targets:'
	@egrep '^(.+)\:\ .*##\ (.+)' ${MAKEFILE_LIST} | sed 's/:.*##/#/' | column -t -c 2 -s '#'

.PHONY: all ## Clean, format, build and test
all: clean-all gofmt build test

.PHONY: install
install:    ## build and install go application executable
	@go install -v ./...

.PHONY: env
env:    ## Print useful environment variables to stdout
	@echo "CURDIR:  $(CURDIR)"
	@echo "SERVICE: $(SERVICE)"
	@echo "PACKAGE: $(PACKAGE)"
	@echo "VERSION: $(VERSION)"

.PHONY: clean
clean:  ## go clean
	@go clean

.PHONY: clean-all
clean-all:  ## remove all generated artifacts and clean all build artifacts
	@go clean -i ./...
	@rm -fr bin

.PHONY: tools
tools:  ## fetch and install all required tools
	@go get -u golang.org/x/tools/cmd/goimports
	@go get -u golang.org/x/lint/golint

.PHONY: fmt
fmt:    ## format the go source files
	@go fmt ./...
	@goimports -w $(FILES)

.PHONY: lint
lint:   ## run go lint on the source files
	@golint $(PACKAGES)

.PHONY: vet
vet:    ## run go vet on the source files
	@go vet ./...

.PHONY: gofmt
gofmt: fmt lint vet

.PHONY: build
build:
	@go build -o $(SERVICE) main.go

.PHONY: test
test:  ## generate grpc code and run short tests
	@go test -v ./... -short

