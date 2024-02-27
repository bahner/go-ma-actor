#!/usr/bin/make -ef

export NAME = go-ma-actor
MODULE_NAME = github.com/bahner/go-ma-actor
export VERSION = "v0.0.2"

GO ?= go
PREFIX ?= /usr/local
KEYSET = go-ma-create-keyset
FETCH = go-ma-fetch-document
ALL =  $(FETCH) $(KEYSET) $(NAME)
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
	
$(NAME): tidy
	$(GO) build -o $(NAME) ./cmd/actor

$(FETCH): tidy
	$(GO) build -o $(FETCH) ./cmd/fetch_document
	
$(KEYSET): tidy
	$(GO) build -o $(KEYSET) ./cmd/create_keyset
	
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
