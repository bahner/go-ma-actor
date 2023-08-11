package main

import (
	"context"
	"flag"

	"github.com/sirupsen/logrus"
	"go.deanishe.net/env"
)

var (
	logLevel   string = env.Get("GO_MYSPACE_LOG_LEVEL", "error")
	rendezvous string = env.Get("GO_MYSPACE_RENDEZVOUS", "myspace")
	nick       string = env.Get("USER", "ghost")
	room       string = env.Get("GO_MYSPACE_ROOM", "mytopic")
)

func Init(ctx context.Context) {

	// Flags - user configurations
	flag.StringVar(&logLevel, "loglevel", logLevel, "Loglevel to use for application")
	flag.StringVar(&rendezvous, "rendezvous", rendezvous, "Unique string to identify group of nodes. Share this with your friends to let them connect with you")
	flag.StringVar(&nick, "nick", nick, "Cosmetic nick to use")
	flag.StringVar(&room, "room", room, "Room to join. This is obviously a TODO as we need more.")

	flag.Parse()

	// Init logger
	log = logrus.New()
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(level)
	log.Info("Logger initialized")

}
