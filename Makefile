#!/usr/bin/make -ef

export NAME = go-ma-actor
MODULE_NAME = github.com/bahner/go-ma-actor
export VERSION = "v0.0.2"

GO ?= go
BUILDFLAGS ?= -ldflags="-s -w"
PREFIX ?= /usr/local
TAR = tar cf -
ZSTD = zstdmt --ultra -22 -o
ZIP = 7z u -tzip
UPX = -qqq
CKSUM = sha256sum
CC = gcc
GOOS = $(shell go env GOOS)
GOARCH = $(shell go env GOARCH)

BINDIR = $(PREFIX)/bin
RELEASEDIR = releases
SIGNATURE_FILE = signature
SIGNATURES = $(RELEASEDIR)/$(SIGNATURE_FILE)

CMDS = actor relay node robot pong document keyset

ANDROID = android-arm64
DARWIN = darwin-amd64 darwin-arm64
FREEBSD = freebsd-amd64 freebsd-arm64
LINUX = linux-amd64 linux-mips64 linux-mips64le linux-ppc64 linux-ppc64le linux-s390x linux-arm64
NETBSD = netbsd-amd64 netbsd-arm64
OPENBSD = openbsd-amd64 openbsd-arm64
WINDOWS =  windows-386 windows-amd64 windows-arm64

POSIX = $(ANDROID) $(DARWIN) $(FREEBSD) $(LINUX) $(NETBSD) $(OPENBSD)

ARM64 = android-arm64 darwin-arm64 netbsd-arm64 openbsd-arm64 
AMD64 = darwin-amd64 freebsd-amd64 linux-amd64 netbsd-amd64 openbsd-amd64 windows-amd64

PLATFORMS =  $(ANDROID) $(DARWIN) $(FREEBSD) $(LINUX) $(NETBSD) $(OPENBSD) $(WINDOWS)

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

ifeq ($(GOOS), windows)
    EXTENSION = .exe
endif

default: clean tidy $(NAME)

all: $(addprefix build-,$(PLATFORMS))

local: clean tidy install

$(BINDIR):
	test -d $(BINDIR) || mkdir -p $(BINDIR)

install: build-$(GOOS)-$(GOARCH) $(BINDIR)
	@$(foreach cmd,$(CMDS), \
		echo Installing $(cmd) for $(GOOS)-$(GOARCH); \
		sudo install -m755 $(GOOS)-$(GOARCH)/$(cmd)$(EXTENSION) $(DESTDIR)$(BINDIR)/; \
	)

debug: BUILDFLAGS += -tags=debug
debug: install

tidy: go.mod
	$(GO) mod tidy

clean:
	rm -rf $(PLATFORMS) $(RELEASEDIR)/* $(CMDS) *.tar *.zip *.log

distclean: clean
	rm -rf $(RELEASEDIR) $(shell git ls-files --others)

release: VERSION := $(shell ./.version)
release: clean $(RELEASEDIR) \
	$(addsuffix .tar,$(RELEASEDIR)/$(POSIX)) \
	$(addsuffix .zip,$(RELEASEDIR)/$(WINDOWS))

	#git tag -a $(VERSION) -m "Release $(VERSION)"

$(SIGNATURES): $(RELEASEDIR)

	cd $(RELEASEDIR) && \
	$(CKSUM) *.zip *.tar > $(SIGNATURE_FILE)

$(RELEASEDIR):
	mkdir -p $(RELEASEDIR)

# Dynamic build commands
build-%: GOOS = $(firstword $(subst -, ,$*))
build-%: GOARCH = $(word 2, $(subst -, ,$*))
build-%: BUILD = $(GOOS)-$(GOARCH)
build-%: $(addprefix $(BUILD),$(CMDS))
	$(eval TARGET_GOOS := $(word 1, $(subst -, ,$*)))
	$(eval TARGET_GOARCH := $(word 2, $(subst -, ,$*)))

    # Conditionally set EXTENSION for Windows targets
	$(if $(findstring windows,$(TARGET_GOOS)), $(eval EXTENSION := .exe))

	@echo "Building for $(TARGET_GOOS)-$(TARGET_GOARCH)..."
	$(foreach cmd,$(CMDS), \
		echo "Building $(cmd) for $(TARGET_GOOS)-$(TARGET_GOARCH)"; \
		mkdir -p $(TARGET_GOOS)-$(TARGET_GOARCH); \
		GOOS=$(TARGET_GOOS) GOARCH=$(TARGET_GOARCH) $(GO) build \
			-o $(TARGET_GOOS)-$(TARGET_GOARCH)/$(cmd)$(EXTENSION) \
			$(BUILDFLAGS) ./cmd/$(cmd) || exit 1; \
	)

# Dynamic release packaging
$(RELEASEDIR)/%: GOOS = $(firstword $(subst -, ,$*))
$(RELEASEDIR)/%: GOARCH = $(firstword $(subst ., ,$(word 2, $(subst -, ,$*))))
$(RELEASEDIR)/%: BUILD = $(GOOS)-$(GOARCH)
$(RELEASEDIR)/%: $(RELEASEDIR) $(BUILD)/$(CMDS)
	@echo "Packaging $(BUILD) build..."
	@if [ "$(GOOS)" = "windows" ]; then \
		$(ZIP) $(RELEASEDIR)/$(NAME)-$(BUILD).zip $(BUILD)/*$(EXTENSION); \
	else \
		$(TAR) $(BUILD) | $(ZSTD) $(RELEASEDIR)/$(BUILD).tar; \
	fi

android-arm64: export CC=aarch64-linux-android-gcc
android-arm64:
	$(MAKE) build-$@
	$(MAKE) $(RELEASEDIR)/$@.tar

darwin-amd64: export CC=x86_64-apple-darwin20-clang
darwin-amd64:
	$(MAKE) build-$@
	$(MAKE) $(RELEASEDIR)/$@.tar

darwin-arm64: export CC=aarch64-apple-darwin20-clang
darwin-arm64:
	$(MAKE) build-$@
	$(MAKE) $(RELEASEDIR)/$@.tar

freebsd-amd64: export CC=x86_64-unknown-freebsd13-clang
freebsd-amd64:
	$(MAKE) build-$@
	$(MAKE) $(RELEASEDIR)/$@.tar

freebsd-arm64: export CC=aarch64-unknown-freebsd13-clang
freebsd-arm64:
	$(MAKE) build-$@
	$(MAKE) $(RELEASEDIR)/$@.tar

linux-amd64: export CC=x86_64-linux-musl-gcc
linux-amd64:
	$(MAKE) build-$@
	$(MAKE) $(RELEASEDIR)/$@.tar

linux-mips64: export CC=mips64-linux-musl-gcc
linux-mips64:
	$(MAKE) build-$@
	$(MAKE) $(RELEASEDIR)/$@.tar

linux-mips64le: export CC=mips64el-linux-musl-gcc
linux-mips64le:
	$(MAKE) build-$@
	$(MAKE) $(RELEASEDIR)/$@

linux-ppc64: export CC=powerpc64-linux-musl-gcc
linux-ppc64:
	$(MAKE) build-$@
	$(MAKE) $(RELEASEDIR)/$@

linux-ppc64le: export CC=powerpc64le-linux-musl-gcc
linux-ppc64le:
	$(MAKE) build-$@
	$(MAKE) $(RELEASEDIR)/$@

linux-s390x: export CC=s390x-linux-musl-gcc
linux-s390x:
	$(MAKE) build-$@
	$(MAKE) $(RELEASEDIR)/$@

linux-arm64: export CC=aarch64-linux-musl-gcc
linux-arm64:
	$(MAKE) build-$@
	$(MAKE) $(RELEASEDIR)/$@

netbsd-amd64: export CC=x86_64-unknown-netbsd9-clang
netbsd-amd64:
	$(MAKE) build-$@
	$(MAKE) $(RELEASEDIR)/$@

netbsd-arm64: export CC=aarch64-unknown-netbsd9-clang
netbsd-arm64:
	$(MAKE) build-$@
	$(MAKE) $(RELEASEDIR)/$@

openbsd-amd64: export CC=x86_64-unknown-openbsd7-clang
openbsd-amd64:
	$(MAKE) build-$@
	$(MAKE) $(RELEASEDIR)/$@

openbsd-arm64: export CC=aarch64-unknown-openbsd7-clang
openbsd-arm64:
	$(MAKE) build-$@
	$(MAKE) $(RELEASEDIR)/$@

windows-386: export CC=i686-w64-mingw32-gcc
windows-386:
	$(MAKE) build-$@
	$(MAKE) $(RELEASEDIR)/$@

windows-amd64: export CC=x86_64-w64-mingw32-gcc
windows-amd64:
	$(MAKE) build-$@
	$(MAKE) $(RELEASEDIR)/$@

windows-arm64: export CC=aarch64-w64-mingw32-gcc
windows-arm64:
	$(MAKE) build-$@
	$(MAKE) $(RELEASEDIR)/$@

android: $(ANDROID)
darwin: $(DARWIN)
freebsd: $(FREEBSD)
linux: $(LINUX)
netbsd: $(NETBSD)
openbsd: $(OPENBSD)
windows: $(WINDOWS)

arm64: $(ARM64)
amd64: $(AMD64)

lint:
	find -name "*.yaml" -exec yamllint -c .yamllintrc {} \;

.PHONY: default all local install debug tidy clean distclean release lint \
	android darwin freebsd linux netbsd openbsd windows \
	arm64 amd64
