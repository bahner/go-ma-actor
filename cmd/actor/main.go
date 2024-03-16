package main

import (
	"context"
	"fmt"
	"os"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/config/db"
	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma-actor/entity/actor"
	"github.com/bahner/go-ma-actor/mode/pong"
	"github.com/bahner/go-ma-actor/p2p"
	"github.com/bahner/go-ma-actor/ui"
	"github.com/spf13/pflag"

	log "github.com/sirupsen/logrus"
)

func main() {

	var (
		err error
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

	// P2P Relay mode
	// Relay mode doesn't need either ui or an actor.
	// So let's just start it quickly and stop here.
	if config.RelayMode() {
		fmt.Println("Starting relay mode...")
		go p2P.StartDiscoveryLoop(context.Background())
		startWebServer(p2P, nil)
		os.Exit(0) // This won't be reached.
	}

	// ACTOR
	a := initActorOrPanic()

	// Start the webserver in the background. Ignore - but log - errors.
	go startWebServer(p2P, a)

	// Pong mode needs the
	if config.PongMode() {
		// In Pong we can just stop here. We dont' need to display anything.
		// or handle input events. Hence this is a blocking call.
		log.Infof("Running in Pong mode")
		pong.Run(a, p2P)
		log.Warnf("Pong run loop ended, exiting...")
		os.Exit(0)

	}

	// We have a valid actor, but for it to be useful, we need to discover peers.
	// discoverPeersOrPanic(p2P)

	ui := initUiOrPanic(p2P, a)

	// START THE ACTOR UI
	fmt.Println("Starting the actor...")
	if err := ui.Run(); err != nil {
		log.Errorf("error running text UI: %s", err)
	}
}

func initActorOrPanic() *actor.Actor {
	// The actor is needed for initialisation of the WebHandler.
	fmt.Println("Creating actor from keyset...")
	a, err := actor.NewFromKeyset(config.ActorKeyset())
	if err != nil {
		log.Debugf("error creating actor: %s", err)
	}

	id := a.Entity.DID.Id

	fmt.Println("Creating and setting DID Document for actor...")
	err = a.CreateAndSetDocument(id)
	if err != nil {
		panic(fmt.Sprintf("error creating document: %s", err))
	}

	// Better safe than sorry.
	// Without a valid actor, we can't do anything.
	if a == nil || a.Verify() != nil {
		panic(fmt.Sprintf("%s is not a valid actor: %v", id, err))
	}

	_, err = entity.GetOrCreateFromDID(a.Entity.DID, false)
	if err != nil {
		panic(fmt.Sprintf("error getting or creating entity: %s", err))
	}

	return a
}

func initUiOrPanic(p2P *p2p.P2P, a *actor.Actor) *ui.ChatUI {
	fmt.Println("Creating text UI...")
	ui, err := ui.New(p2P, a)
	if err != nil {
		panic(fmt.Sprintf("error creating text UI: %s", err))
	}
	return ui
}
