#!/usr/bin/make -f

GO=go
GB=gb

all: build

build: clean test
	@echo "Building..."
	env GOOS=linux GOARCH=amd64 $(GB) build

build-all: build
	@echo "Building others..."
	env GOOS=linux GOARCH=386 $(GB) build
	env GOOS=darwin GOARCH=amd64 $(GB) build
	env GOOS=darwin GOARCH=386 $(GB) build

clean:
	rm -fR pkg bin

test:
	@echo "Running tests..."
	@$(GB) test -test.v=true
