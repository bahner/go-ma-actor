package main

import (
	"context"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/entity/actor"
	"github.com/bahner/go-ma-actor/p2p"
	"github.com/bahner/go-ma-actor/ui/web"
	log "github.com/sirupsen/logrus"
)

const name = "node"

func main() {

	ctx := context.Background()

	// Init config and logger
	config.SetProfile(name)
	actor.InitConfig(config.Profile())

	p, err := p2p.Init(p2p.DefaultOptions())
	if err != nil {
		log.Fatalf("Error initialising P2P: %v", err)
	}

	// Init of actor requires P2P to be initialized
	a := actor.Init()

	go web.Start(web.NewEntityHandler(p, a.Entity))
	go p.StartDiscoveryLoop(ctx)
	go a.Subscribe(ctx, a.Entity)

	// Start application
	StartApplication(p)

	select {}
}
