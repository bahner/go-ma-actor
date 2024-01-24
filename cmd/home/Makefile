#!/usr/bin/make -ef


NAME = go-ma-home
PREFIX ?= /usr/local

default: clean $(NAME)

all: clean $(NAME) install

$(NAME):
	go build -o $(NAME) .

install: $(NAME)
	sudo mv $(NAME) $(PREFIX)/bin

clean:
	rm -f $(NAME)
