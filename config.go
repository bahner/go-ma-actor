package main

import (
	"flag"

	"github.com/sirupsen/logrus"
	"go.deanishe.net/env"
)

var (
	logLevel    string = env.Get("GO_MYSPACE_CLIENT_LOG_LEVEL", "error")
	rendezvous  string = env.Get("GO_MYSPACE_CLIENT_RENDEZVOUS", "myspace")
	serviceName string = env.Get("GO_MYSPACE_CLIENT_SERVICE_NAME", "myspace")
	nick        string = env.Get("USER", "ghost")
	room        string = env.Get("GO_MYSPACE_CLIENT_ROOM", "mytopic")
	secret      string = env.Get("GO_MYSPACE_CLIENT_IDENTITY", "")

	generate *bool
)

func initConfig() {

	// Flags - user configurations
	flag.StringVar(&logLevel, "loglevel", logLevel, "Loglevel to use for application")
	flag.StringVar(&rendezvous, "rendezvous", rendezvous, "Unique string to identify group of nodes. Share this with your friends to let them connect with you")
	flag.StringVar(&serviceName, "servicename", serviceName, "serviceName to use for MDNS discovery")
	flag.StringVar(&nick, "nick", nick, "Cosmetic nick to use")
	flag.StringVar(&room, "room", room, "Room to join. This is obviously a TODO as we need more.")

	generate = flag.Bool("generate", false, "Generate a new private key, prints it and exit the program.")
	flag.StringVar(&secret, "identity", secret, "Base58 encoded secret key used to identofy the client. You.")

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
