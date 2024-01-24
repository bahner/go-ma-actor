#!/usr/bin/make -ef


NAME = go-ma-home
PREFIX ?= /usr/local

default: clean $(NAME)

all: clean $(NAME) install

$(NAME):
	go build -o $(NAME) .

install: $(NAME)
	sudo install -m755 $(NAME) $(DESTDIR)$(PREFIX)/bin/$(NAME)

clean:
	rm -f $(NAME)
