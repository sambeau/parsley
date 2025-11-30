.PHONY: build test clean install version

VERSION := $(shell cat VERSION)
LDFLAGS := -ldflags "-X main.Version=$(VERSION)"

# Build from cmd/pars
build:
	go build $(LDFLAGS) -o pars ./cmd/pars

test:
	go test ./...

clean:
	rm -f pars

install: build
	cp pars $(GOPATH)/bin/

version:
	@echo $(VERSION)

run: build
	./pars

.DEFAULT_GOAL := build
