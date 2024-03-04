#!/usr/bin/make -ef

export NAME = go-ma-actor
MODULE_NAME = github.com/bahner/go-ma-actor
export VERSION = "v0.0.2"

GO ?= go
# This is required for sqlite3 cross-compilation
BUILDFLAGS ?= -ldflags="-s -w"
XZ ?= xz -zf
PREFIX ?= /usr/local
KEYSET = $(NAME)-create-keyset
FETCH = $(NAME)-fetch-document
DEBUG = $(NAME)-debug
ANDROID = android-arm64
DARWIN = darwin-amd64 darwin-arm64
FREEBSD = freebsd-amd64 freebsd-arm64
LINUX = linux-amd64 linux-arm64 linux-mips64 linux-mips64le linux-ppc64 linux-ppc64le linux-s390x
NETBSD = netbsd-amd64 netbsd-arm64
OPENBSD = openbsd-amd64 openbsd-arm64
WINDOWS =  windows-386 windows-amd64
PLATFORMS =  $(ANDROID) $(DARWIN) $(FREEBSD) $(LINUX) $(NETBSD) $(OPENBSD) $(WINDOWS)
ARM64=android-arm64 darwin-arm64 netbsd-arm64 openbsd-arm64 
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
release: clean $(RELEASES) $(PLATFORMS)
	git tag -a $(VERSION) -m "Release $(VERSION)"


$(RELEASES): 
	mkdir -p $(RELEASES)

android: $(ANDROID)

darwin: $(DARWIN)

freebsd: $(FREEBSD)

linux: $(LINUX)

netbsd: $(NETBSD)

openbsd: $(OPENBSD)

windows: $(WINDOWS)

# I need to build these on another computer, so they get their own targets
arm64: $(ARM64)

android-arm64: GOOS=android
android-arm64: GOARCH=arm64
android-arm64: FILENAME = $(RELEASES)/$(NAME)-$(GOOS)-$(GOARCH)
android-arm64: $(RELEASES)
	$(GO) build -o $(FILENAME) $(BUILDFLAGS) ./cmd/actor
	$(XZ) $(FILENAME)

darwin-amd64: GOOS=darwin
darwin-amd64: GOARCH=amd64
darwin-amd64: FILENAME = $(RELEASES)/$(NAME)-$(GOOS)-$(GOARCH)
darwin-amd64: $(RELEASES)
	$(GO) build -o $(FILENAME) $(BUILDFLAGS) ./cmd/actor
	$(XZ) $(FILENAME)

darwin-arm64: GOOS=darwin
darwin-arm64: GOARCH=arm64
darwin-arm64: FILENAME = $(RELEASES)/$(NAME)-$(GOOS)-$(GOARCH)
darwin-arm64: $(RELEASES)
	$(GO) build -o $(FILENAME) $(BUILDFLAGS) ./cmd/actor
	$(XZ) $(FILENAME)

freebsd-amd64: GOOS=freebsd
freebsd-amd64: GOARCH=amd64
freebsd-amd64: FILENAME = $(RELEASES)/$(NAME)-$(GOOS)-$(GOARCH)
freebsd-amd64: $(RELEASES)
	$(GO) build -o $(FILENAME) $(BUILDFLAGS) ./cmd/actor
	$(XZ) $(FILENAME)

freebsd-arm64: GOOS=freebsd
freebsd-arm64: GOARCH=arm64
freebsd-arm64: FILENAME = $(RELEASES)/$(NAME)-$(GOOS)-$(GOARCH)
freebsd-arm64: $(RELEASES)
	$(GO) build -o $(FILENAME) $(BUILDFLAGS) ./cmd/actor
	$(XZ) $(FILENAME)

linux-amd64: GOOS=linux
linux-amd64: GOARCH=amd64
linux-amd64: FILENAME = $(RELEASES)/$(NAME)-$(GOOS)-$(GOARCH)
linux-amd64: $(RELEASES)
	$(GO) build -o $(FILENAME) $(BUILDFLAGS) ./cmd/actor
	$(XZ) $(FILENAME)

linux-arm64: GOOS=linux
linux-arm64: GOARCH=arm64
linux-arm64: CGO_ENABLED=1
linux-arm64: CC=aarch64-linux-musl-gcc
linux-arm64: CGO_CFLAGS="-fPIC"
linux-arm64: CGO_LDFLAGS="-static"
linux-arm64: FILENAME = $(RELEASES)/$(NAME)-$(GOOS)-$(GOARCH)
linux-arm64: $(RELEASES)
	$(GO) build -o $(FILENAME) $(BUILDFLAGS) ./cmd/actor
	$(XZ) $(FILENAME)

linux-mips64: GOOS=linux
linux-mips64: GOARCH=mips64
linux-mips64: FILENAME = $(RELEASES)/$(NAME)-$(GOOS)-$(GOARCH)
linux-mips64: $(RELEASES)
	$(GO) build -o $(FILENAME) $(BUILDFLAGS) ./cmd/actor
	$(XZ) $(FILENAME)

linux-mips64le: GOOS=linux
linux-mips64le: GOARCH=mips64le
linux-mips64le: FILENAME = $(RELEASES)/$(NAME)-$(GOOS)-$(GOARCH)
linux-mips64le: $(RELEASES)
	$(GO) build -o $(FILENAME) $(BUILDFLAGS) ./cmd/actor
	$(XZ) $(FILENAME)

linux-ppc64: GOOS=linux
linux-ppc64: GOARCH=ppc64
linux-ppc64: FILENAME = $(RELEASES)/$(NAME)-$(GOOS)-$(GOARCH)
linux-ppc64: $(RELEASES)
	$(GO) build -o $(FILENAME) $(BUILDFLAGS) ./cmd/actor
	$(XZ) $(FILENAME)

linux-ppc64le: GOOS=linux
linux-ppc64le: GOARCH=ppc64le
linux-ppc64le: FILENAME = $(RELEASES)/$(NAME)-$(GOOS)-$(GOARCH)
linux-ppc64le: $(RELEASES)
	$(GO) build -o $(FILENAME) $(BUILDFLAGS) ./cmd/actor
	$(XZ) $(FILENAME)

linux-s390x: GOOS=linux
linux-s390x: GOARCH=s390x
linux-s390x: FILENAME = $(RELEASES)/$(NAME)-$(GOOS)-$(GOARCH)
linux-s390x: $(RELEASES)
	$(GO) build -o $(FILENAME) $(BUILDFLAGS) ./cmd/actor
	$(XZ) $(FILENAME)

netbsd-amd64: GOOS=netbsd
netbsd-amd64: GOARCH=amd64
netbsd-amd64: FILENAME = $(RELEASES)/$(NAME)-$(GOOS)-$(GOARCH)
netbsd-amd64: $(RELEASES)
	$(GO) build -o $(FILENAME) $(BUILDFLAGS) ./cmd/actor
	$(XZ) $(FILENAME)

netbsd-arm64: GOOS=netbsd
netbsd-arm64: GOARCH=arm64
netbsd-arm64: FILENAME = $(RELEASES)/$(NAME)-$(GOOS)-$(GOARCH)
netbsd-arm64: $(RELEASES)
	$(GO) build -o $(FILENAME) $(BUILDFLAGS) ./cmd/actor
	$(XZ) $(FILENAME)

openbsd-amd64: GOOS=openbsd
openbsd-amd64: GOARCH=amd64
openbsd-amd64: FILENAME = $(RELEASES)/$(NAME)-$(GOOS)-$(GOARCH)
openbsd-amd64: $(RELEASES)
	$(GO) build -o $(FILENAME) $(BUILDFLAGS) ./cmd/actor
	$(XZ) $(FILENAME)

openbsd-arm64: GOOS=openbsd
openbsd-arm64: GOARCH=arm64
openbsd-arm64: FILENAME = $(RELEASES)/$(NAME)-$(GOOS)-$(GOARCH)
openbsd-arm64: $(RELEASES)
	$(GO) build -o $(FILENAME) $(BUILDFLAGS) ./cmd/actor
	$(XZ) $(FILENAME)

windows-386: GOOS=windows
windows-386: GOARCH=386
windows-386: CGO_ENABLED=1
windows-386: CC=i686-w64-mingw32-gcc
windows-386: FILENAME = actor.exe
windows-386: BUILDDIR=$(GOOS)-$(GOARCH)
windows-386: $(RELEASES)
	mkdir -p $(BUILDDIR)
	# $(GO) build -o $(BUILDDIR)/$(FILENAME) $(BUILDFLAGS) ./cmd/actor
	$(GO) build -o $(BUILDDIR)/$(FILENAME) $(BUILDFLAGS) ./cmd/actor
	zip -j $(RELEASES)/$(NAME)-$(GOOS)-$(GOARCH).zip $(BUILDDIR)/$(FILENAME)

windows-amd64: GOOS=windows
windows-amd64: GOARCH=amd64
windows-amd64: CGO_ENABLED=1
windows-amd64: CC=x86_64-w64-mingw32-gcc
windows-amd64: FILENAME = actor.exe
windows-amd64: BUILDDIR=$(GOOS)-$(GOARCH)
windows-amd64: $(RELEASES)
	mkdir -p $(BUILDDIR)
	$(GO) build -o $(BUILDDIR)/$(FILENAME) $(BUILDFLAGS) ./cmd/actor
	zip -j $(RELEASES)/$(NAME)-$(GOOS)-$(GOARCH).zip $(BUILDDIR)/$(FILENAME)

windows-arm64: GOOS=windows
windows-arm64: GOARCH=arm64
windows-arm64: CGO_ENABLED=1
windows-arm64: CC=aarch64-w64-mingw32-gcc
windows-arm64: FILENAME = actor.exe
windows-arm64: BUILDDIR=$(GOOS)-$(GOARCH)
windows-arm64: $(RELEASES)
	mkdir -p $(BUILDDIR)
	$(GO) build -o $(BUILDDIR)/$(FILENAME) $(BUILDFLAGS) ./cmd/actor
	zip -j $(RELEASES)/$(NAME)-$(GOOS)-$(GOARCH).zip $(BUILDDIR)/$(FILENAME)

install: $(BIN)

lint:
	find -name "*.yaml" -exec yamllint -c .yamllintrc {} \;

.PHONY: default init tidy build client serve install clean distclean lint
