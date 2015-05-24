#!/usr/bin/make -f

SHELL=/bin/bash
MKDIR=mkdir
GIT=git
GO=go
RM=rm -rf
CROSS=https://github.com/davecheney/golang-crosscompile.git
CROSS_BASH=source golang-crosscompile/crosscompile.bash

SOURCE=./...
TARGETS=darwin-386 darwin-amd64 linux-386 linux-amd64 linux-arm

all: test

build: deps
	@echo "Building close..."
	@$(GO) build -o bin/close $(SOURCE)

deps:
	@echo "Downloading libraries..."
	@$(GO) get gopkg.in/alecthomas/kingpin.v1
	@$(GO) get github.com/Sirupsen/logrus
	@$(GO) get github.com/google/go-github/github
	@$(GO) get github.com/stretchr/testify/assert

golang-crosscompile:
	$(GIT) clone $(CROSS)
	$(CROSS_BASH) && \
	go-crosscompile-build-all

xbuild: deps golang-crosscompile dirs
	@for target in $(TARGETS); do \
		echo "Building close for $$target..."; \
		$(CROSS_BASH) && \
		$(GO)-$$target build -o bin/close-$$target $(SOURCE); \
	done;

dirs:
	@$(MKDIR) -p bin

test: build
	@echo "Run tests..."
	@$(GO) test ./...

clean:
	@echo "Cleanup binaries..."
	$(RM) bin

realclean: clean
	$(RM) golang-crosscompile

coverage:
	# This is a script provided by upstream. We won't to need this in 1.5 of Golang.
	scripts/coverage.sh