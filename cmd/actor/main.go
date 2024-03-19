package main

import (
	"fmt"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/config/db"
	"github.com/bahner/go-ma-actor/entity/actor"
	"github.com/bahner/go-ma-actor/p2p"
	"github.com/bahner/go-ma-actor/ui"
	"github.com/bahner/go-ma-actor/ui/web"

	log "github.com/sirupsen/logrus"
)

func main() {

	var (
		err error
	)

	fmt.Println("Initialising actor configuation...")
	actor.InitConfig(config.Profile())

	// DB
	fmt.Println("Initialising DB ...")
	_, err = db.Init()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize db: %v", err))
	}

	// P2P
	fmt.Println("Setting default p2p options...")
	p2pOpts := p2p.DefaultOptions()
	fmt.Println("Initialising p2p...")
	p2P, err := p2p.Init(p2pOpts)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize p2p: %v", err))
	}

	// ACTOR
	fmt.Println("Initialising actor...")
	a := actor.Init()

	// WEB
	fmt.Println("Initialising web UI...")
	wh := web.NewEntityHandler(p2P, a.Entity)
	go web.Start(wh)

	// TEXT UI
	fmt.Println("Initialising text UI...")
	ui := ui.Init(p2P, a)

	// START THE ACTOR UI
	fmt.Println("Starting the actor...")
	if err := ui.Run(); err != nil {
		log.Errorf("error running text UI: %s", err)
	}
}
