package main

import (
	"fmt"

	"github.com/bahner/go-ma-actor/config/db"
	"github.com/bahner/go-ma-actor/entity/actor"
	"github.com/bahner/go-ma-actor/ui"
	"github.com/bahner/go-ma-actor/ui/web"

	log "github.com/sirupsen/logrus"
)

func main() {

	var (
		err error
	)

	actor.InitConfig()

	// DB
	fmt.Println("Initialising DB ...")
	_, err = db.Init()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize db: %v", err))
	}

	// P2P
	p2P, err := initP2P()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize p2p: %v", err))
	}

	// ACTOR
	a := actor.Init()

	// Start the webserver in the background. Ignore - but log - errors.
	wh := web.NewWebEntityHandler(p2P, a.Entity)
	go web.Start(wh)

	// We have a valid actor, but for it to be useful, we need to discover peers.
	// discoverPeersOrPanic(p2P)

	ui := ui.Init(p2P, a)

	// START THE ACTOR UI
	fmt.Println("Starting the actor...")
	if err := ui.Run(); err != nil {
		log.Errorf("error running text UI: %s", err)
	}
}
