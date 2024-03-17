package main

import (
	"fmt"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/config/db"
	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma-actor/entity/actor"
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
	config.InitCommonFlags()
	config.InitActorFlags()
	pflag.Parse()
	config.Init()

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
