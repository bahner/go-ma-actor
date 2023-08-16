#!/usr/bin/make -ef

NAME = go-ma-actor
MODULE_NAME = github.com/bahner/go-ma-actor

GO ?= go
PREFIX ?= /usr/local

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

default: clean tidy $(NAME)

init: go.mod tidy

go.mod:
	$(GO) mod init $(MODULE_NAME)

tidy: go.mod
	$(GO) mod tidy

$(NAME): tidy
	$(GO) build -o $(NAME)

clean:
	rm -f $(NAME)

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

install:
	install -Dm755 $(NAME) $(DESTDIR)$(PREFIX)/bin/$(NAME)

run: clean $(NAME)
	./$(NAME)

up:
	docker-compose up -d --remove-orphans

vault:
	docker-compose up -d vault

.PHONY: default init tidy build client serve install clean distclean
