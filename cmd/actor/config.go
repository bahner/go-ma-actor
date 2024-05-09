package main

import (
	"os"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma-actor/entity/actor"
	"github.com/bahner/go-ma-actor/p2p/peer"

	log "github.com/sirupsen/logrus"
)

func initConfig(defaultProfileName string) actor.ActorConfig {

	config.SetDefaultProfileName(defaultProfileName)
	c := actor.Config()

	if config.GenerateFlag() {
		config.GenerateConfig(&c)
	}

	if config.ShowConfigFlag() {
		c.Print()
	}

	if config.ShowConfigFlag() || config.GenerateFlag() {
		os.Exit(0)
	}

	log.Info("Reading CSV files...")
	go entity.WatchCSV()
	go peer.WatchCSV()

	return c
}
