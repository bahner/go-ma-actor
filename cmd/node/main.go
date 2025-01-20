package main

import (
	"context"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/entity/actor"
	"github.com/bahner/go-ma-actor/p2p"
	"github.com/bahner/go-ma-actor/ui/web"
)

func main() {

	ctx := context.Background()

	Config(config.Profile())

	p2pOpts := p2p.DefaultP2POptions()

	// Init of actor requires P2P to be initialized
	a := actor.Init(p2pOpts)

	go web.Start(web.NewEntityHandler(p, a.Entity))
	go p.StartDiscoveryLoop(ctx)
	go a.Subscribe(ctx, a.Entity)

	// Start application
	StartApplication(p)

	select {}
}
