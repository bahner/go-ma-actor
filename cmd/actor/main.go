package main

import (
	"fmt"

	"github.com/bahner/go-ma-actor/actor"
	"github.com/bahner/go-ma-actor/p2p"
	"github.com/bahner/go-ma-actor/ui"
	"github.com/bahner/go-ma-actor/ui/web"

	log "github.com/sirupsen/logrus"
)

const defaultProfileName = "actor"

func main() {

	initConfig(defaultProfileName)

	// ACTOR
	fmt.Println("Initialising actor...")
	a := actor.Init(p2p.DefaultP2POptions())

	// WEB
	fmt.Println("Initialising web UI...")
	wh := web.NewEntityHandler(a.P2P, a.Entity)
	go web.Start(wh)

	// TEXT UI
	fmt.Println("Initialising text UI...")
	ui := ui.Init(a)

	// START THE ACTOR UI
	fmt.Println("Starting the actor...")
	if err := ui.Run(); err != nil {
		log.Errorf("error running text UI: %s", err)
	}
}
