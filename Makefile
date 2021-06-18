SHELL=/usr/bin/env bash

GOFLAGS+=-ldflags=-X="bee-agent/build.CurrentCommit"="+git$(subst -,.,$(shell git describe --always --match=NeVeRmAtCh --dirty 2>/dev/null || git rev-parse --short HEAD 2>/dev/null))"

.PHONY: default
default:  linux;
all: linux windows darwin

linux:
	rm -f bee-agent-linux-amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bee-agent-linux-amd64 $(GOFLAGS)

windows:
	rm -f bee-agent.exe
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(GOFLAGS)

darwin:
	rm -f bee-agent-darwin-amd64
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bee-agent-darwin-amd64 $(GOFLAGS)