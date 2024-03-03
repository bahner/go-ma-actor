#!/usr/bin/make -ef

export NAME = go-ma-actor
MODULE_NAME = github.com/bahner/go-ma-actor
export VERSION = "v0.0.2"

GO ?= go
BUILDFLAGS ?= -ldflags="-s -w"
TAR ?= tar cJf
PREFIX ?= /usr/local
KEYSET = $(NAME)-create-keyset
FETCH = $(NAME)-fetch-document
DEBUG = $(NAME)-debug
PLATFORMS = linux-amd64 windows-amd64 windows-386 darwin-amd64 darwin-arm64
ALL =  $(FETCH) $(KEYSET) $(NAME) $(DEBUG)
BIN = $(PREFIX)/bin
RELEASES = releases

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

default: clean tidy $(NAME)

all: tidy releases $(PLATFORMS)

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
	rm -rf $(PLATFORMS)
	# rm -f $(NAME)-*.tar
	# find -type f -name "*.log" -delete
	# rm -f actor.exe

distclean: clean
	rm -rf releases
	rm -f $(shell git ls-files --others)


release: VERSION = $(shell ./.version)
release: clean $(RELEASES) windows darwin linux-amd64
	git tag -a $(VERSION) -m "Release $(VERSION)"


$(RELEASES): 
	mkdir -p $(RELEASES)

linux-amd64: $(ALL)
	mkdir -p linux-amd64
	mv $(ALL) linux-amd64
	$(TAR) $(RELEASES)/$(NAME)-linux-amd64.tar -C linux-amd64 $(ALL)

windows-amd64: GOOS=windows
windows-amd64: GOARCH=amd64
windows-amd64: FILENAME = actor.exe
windows-amd64: BUILDDIR=$(GOOS)-$(GOARCH)
windows-amd64: $(RELEASES)
	mkdir -p $(BUILDDIR)
	$(GO) build -o $(BUILDDIR)/$(FILENAME) $(BUILDFLAGS) ./cmd/actor
	zip -j $(RELEASES)/$(NAME)-$(GOOS)-$(GOARCH).zip $(BUILDDIR)/$(FILENAME)

windows-386: GOOS=windows
windows-386: GOARCH=386
windows-386: FILENAME = actor.exe
windows-386: BUILDDIR=$(GOOS)-$(GOARCH)
windows-386: $(RELEASES)
	mkdir -p $(BUILDDIR)
	$(GO) build -o $(BUILDDIR)/$(FILENAME) $(BUILDFLAGS) ./cmd/actor
	zip -j $(RELEASES)/$(NAME)-$(GOOS)-$(GOARCH).zip $(BUILDDIR)/$(FILENAME)



darwin-amd64: GOOS=darwin
darwin-amd64: GOARCH=amd64
darwin-amd64: BUILDDIR=$(GOOS)-$(GOARCH)
darwin-amd64: $(RELEASES)
	mkdir -p $(BUILDDIR)
	$(GO) build -o $(BUILDDIR)/$(NAME) $(BUILDFLAGS) ./cmd/actor
	$(TAR) $(RELEASES)/$(NAME)-$(GOOS)-$(GOARCH).tar -C $(BUILDDIR) $(NAME)	

darwin-arm64: GOOS=darwin
darwin-arm64: GOARCH=arm64
darwin-arm64: BUILDDIR=$(GOOS)-$(GOARCH)
darwin-arm64: $(RELEASES)
	mkdir -p $(BUILDDIR)
	$(GO) build -o $(BUILDDIR)/$(NAME) $(BUILDFLAGS) ./cmd/actor
	$(TAR) $(RELEASES)/$(NAME)-$(GOOS)-$(GOARCH).tar -C $(BUILDDIR) $(NAME)	

install: $(BIN)

lint:
	find -name "*.yaml" -exec yamllint -c .yamllintrc {} \;

.PHONY: default init tidy build client serve install clean distclean lint
