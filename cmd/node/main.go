package main

import (
	"context"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/entity/actor"
	"github.com/bahner/go-ma-actor/p2p"
	"github.com/bahner/go-ma-actor/ui/web"
	log "github.com/sirupsen/logrus"
)

func main() {

	ctx := context.Background()

	Config(config.Profile())

	// Init of actor requires P2P to be initialized
	a := actor.Init()

	p, err := p2p.Init(a.Keyset.Identity, p2p.DefaultOptions())
	if err != nil {
		log.Fatalf("Error initialising P2P: %v", err)
	}

	go web.Start(web.NewEntityHandler(p, a.Entity))
	go p.StartDiscoveryLoop(ctx)
	go a.Subscribe(ctx, a.Entity)

	// Start application
	StartApplication(p)

	select {}
}
