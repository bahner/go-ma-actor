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
ALL =  $(FETCH) $(KEYSET) $(NAME) $(DEBUG)
BIN = $(PREFIX)/bin

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

console:
	docker-compose up -d
	docker attach go-space-pubsub_space_1

distclean: clean
	rm -f $(shell git ls-files --exclude-standard --others)

down:
	docker-compose down

image:
	docker build \
		-t $(IMAGE) \
		--build-arg "BUILD_IMAGE=$(BUILD_IMAGE)" \
		.

install: $(BIN)

lint:
	find -name "*.yaml" -exec yamllint -c .yamllintrc {} \;

webui:

run: clean $(NAME)
	./$(NAME)

up:
	docker-compose up -d --remove-orphans

vault:
	docker-compose up -d vault

.PHONY: default init tidy build client serve install clean distclean
