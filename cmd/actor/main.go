package main

import (
	"fmt"

	"github.com/bahner/go-ma-actor/config/db"

	log "github.com/sirupsen/logrus"
)

func main() {

	var (
		err error
	)

	initConfig()

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

	// PEER
	fmt.Println("Initialising peer ...")
	err = initPeer(p2P.Host.ID().String())
	if err != nil {
		panic(fmt.Sprintf("failed to initialize peer: %v", err))
	}

	// ACTOR
	a := initActorOrPanic()

	// Start the webserver in the background. Ignore - but log - errors.
	go startWebServer(p2P, a)

	// We have a valid actor, but for it to be useful, we need to discover peers.
	// discoverPeersOrPanic(p2P)

	ui := initUiOrPanic(p2P, a)

	// START THE ACTOR UI
	fmt.Println("Starting the actor...")
	if err := ui.Run(); err != nil {
		log.Errorf("error running text UI: %s", err)
	}
}
