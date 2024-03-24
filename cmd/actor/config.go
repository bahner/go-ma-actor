package main

import (
	"os"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/entity/actor"
)

func initConfig(defaultProfileName string) actor.ActorConfig {

	config.SetDefaultProfileName(defaultProfileName)
	c := actor.Config()

	if config.GenerateFlag() {
		config.Generate(&c)
	}

	if config.ShowConfigFlag() {
		c.Print()
	}

	if config.ShowConfigFlag() || config.GenerateFlag() {
		os.Exit(0)
	}

	return c
}
