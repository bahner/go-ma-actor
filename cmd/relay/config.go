package main

import (
	"os"

	"github.com/bahner/go-ma-actor/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

func initConfig(profile string) {

	// Always parse the flags first
	config.InitCommonFlags()
	pflag.Parse()
	config.SetProfile(profile)
	config.Init()

	if config.GenerateFlag() {
		// Reinit logging to STDOUT
		log.SetOutput(os.Stdout)
		log.Info("Generating new actor and node identity")
		node, err := config.GenerateNodeIdentity()
		if err != nil {
			log.Fatalf("Failed to generate node identity: %v", err)
		}
		relayConfig := configTemplate(node)
		config.Generate(relayConfig)
		os.Exit(0)
	}
}
