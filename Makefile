GO ?= go
VERSION ?= git.$(shell git rev-parse --short HEAD)
GOFLAGS := -ldflags "-w -s -X main.Version=${VERSION}"

DEPS_GO_COMMON := ./go.mod ./go.sum $(shell echo ./*.go)
DEPS_PROXYSSH != echo ./proxyssh/*.go
DEPS_BOX != echo ./box/*.go
DEPS_UPPY != echo ./uppy/*.go

all: bin/proxyssh bin/box bin/uppy

bin/proxyssh: $(DEPS_GO_COMMON) $(DEPS_PROXYSSH)
	$(GO) build $(GOFLAGS) -o $@ ./proxyssh/

bin/box: $(DEPS_GO_COMMON) $(DEPS_BOX)
	$(GO) build $(GOFLAGS) -o $@ ./box/

bin/uppy: $(DEPS_GO_COMMON) $(DEPS_UPPY)
	$(GO) build $(GOFLAGS) -o $@ ./uppy/

clean:
	rm -rf bin/

.PHONY: all clean
