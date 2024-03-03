package main

import (
	"fmt"
	"os"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/config/db"
	"github.com/bahner/go-ma-actor/entity/actor"
	"github.com/bahner/go-ma-actor/mode/relay"
	"github.com/bahner/go-ma-actor/p2p"
	"github.com/bahner/go-ma-actor/ui"
	"github.com/spf13/pflag"

	log "github.com/sirupsen/logrus"
)

func main() {

	var (
		err error
		p   *p2p.P2P
	)

	// Always parse the flags first
	pflag.Parse()

	// MODE

	// Then init the config
	// There's a lot of stuff going on in here.
	mode := config.Mode()
	config.Init(mode)
	fmt.Printf("Starting in %s mode.\n", mode)

	// DB
	fmt.Print("Initialising DB ...")
	_, err = db.Init()
	if err != nil {
		log.Fatalf("failed to initialize db: %v", err)
	}
	fmt.Println("done.")
	// P2P

	// Configure and start P2P
	fmt.Print("Initialising libp2p...")

	if config.RelayMode() {
		fmt.Print("Relay mode enabled.")
		p, err = p2p.Init(nil, relay.GetOptions()...)
	} else {
		p, err = p2p.Init(nil)
	}
	if err != nil {
		log.Fatalf("failed to initialize p2p: %v", err)
	}
	fmt.Println("done.")

	if config.RelayMode() {
		fmt.Println("Starting relay mode...")
		startWebServer(p, nil)
		os.Exit(0) // This won't be reached.
	}

	// ACTOR

	// The actor is needed for initialisation of the WebHandler.
	fmt.Print("Creating actor from keyset...")
	a, err := actor.NewFromKeyset(config.ActorKeyset())
	if err != nil {
		log.Debugf("error creating actor: %s", err)
	}
	fmt.Println("done.")

	id := a.Entity.DID.Id

	fmt.Print("Creating and setting DID Document for actor...")
	err = a.CreateAndSetDocument(id)
	if err != nil {
		log.Fatalf("error creating document: %s", err)
	}
	fmt.Println("done.")

	// Better safe than sorry.
	// Without a valid actor, we can't do anything.
	if a == nil || a.Verify() != nil {
		log.Fatalf("%s is not a valid actor: %v", id, err)
	}

	// PEER DISCOVERY

	// We need to discover peers before we can do anything else.
	// So this is a blocking call.
	fmt.Print("Discovering peers...")
	p.DiscoverPeers()
	if err != nil {
		log.Fatalf("failed to initialize p2p: %v", err)
	}
	fmt.Println("done.")

	// UI

	// Draw the UI.
	fmt.Print("Creating text UI...")
	ui, err := ui.NewChatUI(p, a)
	if err != nil {
		log.Fatalf("error creating text UI: %s", err)
	}
	fmt.Println("done.")

	// WEBSERVER

	fmt.Print("Starting web server...")
	go startWebServer(p, a)
	fmt.Println("done.")

	fmt.Println("Starting the actor...")
	if err := ui.Run(); err != nil {
		log.Errorf("error running text UI: %s", err)
	}
}
