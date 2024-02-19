package main

import (
	"context"
	"net/http"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma-actor/entity/actor"
	"github.com/bahner/go-ma-actor/p2p"
	"github.com/bahner/go-ma-actor/ui"
	"github.com/spf13/pflag"

	log "github.com/sirupsen/logrus"
)

func main() {

	pflag.Parse()

	config.Init("actor")
	config.InitLogging()
	config.InitP2P()
	config.InitActor()

	p, err := p2p.Init(nil)
	if err != nil {
		log.Fatalf("failed to initialize p2p: %v", err)
	}

	// We need to discover peers before we can do anything else.
	// So this is a blocking call.
	p.DiscoverPeers()
	if err != nil {
		log.Fatalf("failed to initialize p2p: %v", err)
	}

	// Now we can start continuous discovery in the background.
	go p.DiscoveryLoop(context.Background())

	// The actor is needed for the WebHandler.
	a, err := actor.NewFromKeyset(config.GetKeyset())
	if err != nil {
		log.Debugf("error creating actor: %s", err)
	}

	id := a.Entity.DID.Id

	err = a.CreateAndSetDocument(id)
	if err != nil {
		log.Fatalf("error creating document: %s", err)
	}

	// Better safe than sorry.
	// Without a valid actor, we can't do anything.
	if a == nil || a.Verify() != nil {
		log.Fatalf("%s is not a valid actor: %v", id, err)
	}

	// Start a simple web server to handle incoming requests.
	// This is defined in web.go. It makes it possible to add extra parameters to the handler.
	h := &entity.WebHandlerData{
		P2P:    p,
		Entity: a.Entity,
	}
	http.HandleFunc("/", h.WebHandler)
	log.Infof("Listening on %s", config.GetHttpSocket())
	go http.ListenAndServe(config.GetHttpSocket(), nil)

	// Draw the UI.
	log.Debugf("Starting text UI")
	ui, err := ui.NewChatUI(p, a)
	if err != nil {
		log.Fatalf("error creating text UI: %s", err)
	}

	if err := ui.Run(); err != nil {
		log.Errorf("error running text UI: %s", err)
	}
}
