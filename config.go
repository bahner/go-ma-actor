package main

import (
	"flag"

	log "github.com/sirupsen/logrus"
	"go.deanishe.net/env"
)

var (
	logLevel    string = env.Get("GO_MA_ACTOR_LOG_LEVEL", "error")
	rendezvous  string = env.Get("GO_MA_ACTOR_RENDEZVOUS", "/ma/0.0.1")
	serviceName string = env.Get("GO_MA_ACTOR_SERVICE_NAME", "/ma/0.0.1")
	nick        string = env.Get("USER", "ghost")
	topic       string = env.Get("GO_MA_ACTOR_ROOM", "mytopic")
	keyset      string = env.Get("GO_MA_ACTOR_IDENTITY", "")

	generate *bool
	genenv   *bool
	publish  *bool
)

func initConfig() {

	// Flags - user configurations
	flag.StringVar(&logLevel, "loglevel", logLevel, "Loglevel to use for application")
	flag.StringVar(&rendezvous, "rendezvous", rendezvous, "Unique string to identify group of nodes. Share this with your friends to let them connect with you")
	flag.StringVar(&serviceName, "servicename", serviceName, "serviceName to use for MDNS discovery")
	flag.StringVar(&nick, "nick", nick, "Cosmetic nick to use")
	flag.StringVar(&topic, "topic", topic, "Room (topic) to join. This is obviously a TODO as we need more.")

	// The secret sauce. Use or generate a new one.
	flag.StringVar(&keyset, "identity", keyset, "Base58 encoded secret key used to identify the client. You.")

	generate = flag.Bool("generate", false, "Generate a new private key, prints it and exit the program.")
	genenv = flag.Bool("genenv", false, "Generates a new environment file with a new private key to stdout")
	publish = flag.Bool("publish", false, "Publishes keyset to IPFS when using genenv or generate")

	flag.Parse()

	if *generate || *genenv {
		generateKeyset(nick, *generate)
	}

	// Init logger
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(level)
	log.Info("Logger initialized")

}
