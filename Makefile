#!/usr/bin/make -ef

export NAME = go-ma-actor
MODULE_NAME = github.com/bahner/go-ma-actor
export VERSION = "v0.0.2"

GO ?= go
BUILDFLAGS ?= -ldflags="-s -w"
PREFIX ?= /usr/local
KEYSET = $(NAME)-create-keyset
FETCH = $(NAME)-fetch-document
DEBUG = $(NAME)-debug
PLATFORMS = linux-amd64 windows darwin
ALL =  $(FETCH) $(KEYSET) $(NAME) $(DEBUG)
BIN = $(PREFIX)/bin
RELEASES = releases

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

default: clean tidy $(NAME)

all: tidy $(ALL)

$(BIN): $(ALL)
	test -d $(BIN)
	sudo install -m755 $(ALL) $(DESTDIR)$(BIN)
	
$(DEBUG): BUILDFLAGS = -tags=debug
$(DEBUG): tidy
	$(GO) build -o $(DEBUG) $(BUILDFLAGS) ./cmd/actor

$(NAME): tidy
	$(GO) build -o $(NAME) $(BUILDFLAGS) ./cmd/actor

$(FETCH): tidy
	$(GO) build -o $(FETCH) $(BUILDFLAGS) ./cmd/fetch_document
	
$(KEYSET): tidy
	$(GO) build -o $(KEYSET) $(BUILDFLAGS) ./cmd/create_keyset
	
init: go.mod tidy

go.mod:
	$(GO) mod init $(MODULE_NAME)

tidy: go.mod
	$(GO) mod tidy

clean:
	rm -f $(ALL) 
	rm -rf $(PLATFORMS)
	rm -rf release

distclean: clean
	rm -f $(shell git ls-files --exclude-standard --others)


release: clean $(RELEASES) $(PLATFORMS)

$(RELEASES): 
	mkdir -p $(RELEASES)

linux-amd64: $(ALL)	
	tar cJf $(RELEASES)/$(NAME)-linux-amd64.tar $(ALL)

windows: GOOS=windows
windows: windows-amd64 windows-386

windows-amd64: GOARCH=amd64
windows-amd64: BUILDDIR=$(GOOS)-$(GOARCH)
windows-amd64: $(RELEASES)
	go build -o $(BUILDDIR)/actor.exe ./cmd/actor
	zip $(RELEASES)/$(NAME)-$(GOOS)-$(GOARCH).zip $(BUILDDIR)/actor.exe

windows-386: GOARCH=386
windows-386: BUILDDIR=$(GOOS)-$(GOARCH)
windows-386: $(RELEASES)
	go build -o $(BUILDDIR)/actor.exe ./cmd/actor
	zip $(RELEASES)/$(NAME)-$(GOOS)-$(GOARCH).zip $(BUILDDIR)/actor.exe

darwin: GOOS=darwin
darwin: darwin-amd64 darwin-arm64

darwin-amd64: GOARCH=amd64
darwin-amd64: BUILDDIR=$(GOOS)-$(GOARCH)
darwin-amd64: $(RELEASES)
	mkdir -p $(BUILDDIR)
	go build -o $(BUILDDIR)/$(NAME) ./cmd/actor
	tar cJf $(NAME)-$(GOOS)-$(GOARCH).tar -C $(BUILDDIR) $(NAME)	

darwin-arm64: GOARCH=arm64
darwin-arm64: BUILDDIR=$(GOOS)-$(GOARCH)
darwin-arm64: $(RELEASES)
	mkdir -p $(BUILDDIR)
	go build -o $(BUILDDIR)/$(NAME) ./cmd/actor
	tar cJf $(NAME)-$(GOOS)-$(GOARCH).tar -C $(BUILDDIR) $(NAME)	

install: $(BIN)

lint:
	find -name "*.yaml" -exec yamllint -c .yamllintrc {} \;

.PHONY: default init tidy build client serve install clean distclean $(PLATFORMS) $(RELEASES) $(BIN) $(ALL) $(NAME) $(FETCH) $(KEYSET) $(DEBUG) lint
