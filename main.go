package main

import (
	"flag"
	"os"

	"github.com/bahner/go-ma-actor/actor"
	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/p2p"
	"github.com/bahner/go-ma-actor/ui"

	log "github.com/sirupsen/logrus"
)

func main() {

	flag.Parse()
	config.Init()

	p, err := p2p.Init(nil)
	if err != nil {
		log.Errorf("failed to initialize p2p: %v", err)
		os.Exit(75)
	}

	p.DHT.DiscoverPeers()

	if err != nil {
		log.Errorf("failed to initialize p2p: %v", err)
		os.Exit(75)
	}

	a, err := actor.NewFromKeyset(config.GetKeyset(), config.GetForcePublish())
	if err != nil || a == nil {
		log.Errorf("failed to create actor: %v", err)
		os.Exit(70)
	}

	e := config.GetEntity()
	// Draw the UI.
	log.Debugf("Starting text UI")
	ui := ui.NewChatUI(p, a, e)
	if err := ui.Run(); err != nil {
		log.Errorf("error running text UI: %s", err)
	}
}
