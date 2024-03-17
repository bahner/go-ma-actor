package main

import (
	"os"

	"github.com/bahner/go-ma-actor/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

func initConfig() {

	// Always parse the flags first
	config.InitCommonFlags()
	config.InitActorFlags()
	pflag.Parse()
	config.Init()

	if config.GenerateFlag() {
		// Reinit logging to STDOUT
		log.SetOutput(os.Stdout)
		log.Info("Generating new actor and node identity")
		actor, node := config.GenerateActorIdentitiesOrPanic()
		actorConfig := configTemplate(actor, node)
		config.Generate(actorConfig)
		os.Exit(0)
	}
}
