package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"

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

	fmt.Print("Initialising libp2p...")
	p, err := p2p.Init(nil)
	if err != nil {
		log.Fatalf("failed to initialize p2p: %v", err)
	}
	fmt.Println("done.")

	// Now we can start continuous discovery in the background.
	fmt.Print("Starting discovery loop...")
	go p.DiscoveryLoop(context.Background())
	fmt.Println("done.")

	// The actor is needed for the WebHandler.
	fmt.Print("Creating actor from keyset...")
	a, err := actor.NewFromKeyset(config.GetKeyset())
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

	// We need to discover peers before we can do anything else.
	// So this is a blocking call.
	fmt.Print("Discovering peers...")
	p.DiscoverPeers()
	if err != nil {
		log.Fatalf("failed to initialize p2p: %v", err)
	}
	fmt.Println("done.")

	fmt.Print("Starting web server...")
	// Start a simple web server to handle incoming requests.
	// This is defined in web.go. It makes it possible to add extra parameters to the handler.
	h := &entity.WebHandlerData{
		P2P:    p,
		Entity: a.Entity,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.WebHandler)
	mux.HandleFunc("/profile", pprof.Profile)
	go http.ListenAndServe(config.GetHttpSocket(), mux)
	log.Infof("Listening on %s", config.GetHttpSocket())
	fmt.Println("done.")
	fmt.Println("Web server started on http://" + config.GetHttpSocket() + "/")

	// Draw the UI.
	fmt.Print("Creating text UI...")
	ui, err := ui.NewChatUI(p, a)
	if err != nil {
		log.Fatalf("error creating text UI: %s", err)
	}
	fmt.Println("done.")

	// Run the UI.
	fmt.Println("Starting the actor...")
	if err := ui.Run(); err != nil {
		log.Errorf("error running text UI: %s", err)
	}
}
