package main

import (
	"sync"

	"github.com/ayush6624/go-chatgpt"
	log "github.com/sirupsen/logrus"
)

var (
	once          sync.Once
	chatgptClient *chatgpt.Client
)

func client() *chatgpt.Client {

	var err error

	once.Do(func() {

		chatgptClient, err = chatgpt.NewClient(openAIKey())
		if err != nil {
			log.Fatal(err)
		}

	})

	return chatgptClient

}
