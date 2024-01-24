#!/usr/bin/make -ef

export NAME = go-ma-actor
MODULE_NAME = github.com/bahner/go-ma-actor
export VERSION = "v0.0.2"

GO ?= go
PREFIX ?= /usr/local

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

default: clean tidy $(NAME)

all: clean tidy install

build: $(NAME)

init: go.mod tidy

go.mod:
	$(GO) mod init $(MODULE_NAME)

tidy: go.mod
	$(GO) mod tidy

$(NAME): tidy
	$(GO) build -o $(NAME)

clean:
	rm -f $(NAME)
	make -C cmd/home clean
	make -C cmd/relay clean

console:
	docker-compose up -d
	docker attach go-space-pubsub_space_1

distclean: clean
	rm -f $(shell git ls-files --exclude-standard --others)

down:
	docker-compose down

home:
	make -C cmd/home

relay:
	make -C cmd/relay go-ma-relay

image:
	docker build \
		-t $(IMAGE) \
		--build-arg "BUILD_IMAGE=$(BUILD_IMAGE)" \
		.

install: relay home $(NAME)
	sudo make -C cmd/home install
	sudo make -C cmd/relay install
	sudo install -m755 $(NAME) $(DESTDIR)$(PREFIX)/bin/$(NAME)


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
