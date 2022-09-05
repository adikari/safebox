# TODO: remove hard coded version
VERSION := 1
ifndef VERSION
	VERSION := $(shell git describe --tags --always --dirty="-dev")
endif

ifndef TARGETARCH
	TARGETARCH := $(shell arch)
endif

LDFLAGS := -ldflags='-X "main.Version=$(VERSION)"'

test:
	go test -v ./...

all: darwin-amd64 darwin-arm64 linux-amd64 windows-amd64.exe

clean: |
	rm -rf ./dist
	rm -f ./safebox

dist/:
	mkdir -p dist

build: clean safebox

safebox:
	CGO_ENABLED=0 go build -trimpath $(LDFLAGS) -o $@

linux: linux-$(TARGETARCH)
	cp $^ safebox

darwin-amd64: | dist/
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -trimpath $(LDFLAGS) -o dist/safebox-$(VERSION)-$@

darwin-arm64: | dist/
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -trimpath $(LDFLAGS) -o dist/safebox-$(VERSION)-$@

linux-amd64: | dist/
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath $(LDFLAGS) -o dist/safebox-$(VERSION)-$@

linux-arm64 linux-aarch64: | dist/
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -trimpath $(LDFLAGS) -o dist/safebox-$(VERSION)-$@

windows-amd64.exe: | dist/
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath $(LDFLAGS) -o dist/safebox-$(VERSION)-$@

.PHONY: clean all linux
