package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/pflag"

	"github.com/bahner/go-ma-actor/actor"
	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/p2p"

	log "github.com/sirupsen/logrus"
)

func main() {

	pflag.Parse()

	config.Init("pong")
	config.InitLogging()
	config.InitP2P()
	config.InitActor()

	ctx := context.Background()

	p, err := p2p.Init(nil)
	if err != nil {
		log.Errorf("Error initializing p2p node: %v", err)
		os.Exit(69) // EX_UNAVAILABLE
	}

	if err != nil {
		log.Errorf("Error initializing p2p node: %v", err)
		os.Exit(69) // EX_UNAVAILABLE
	}

	// We need to discover peers before we can do anything else.
	p.DiscoverPeers()

	n := p.Node

	a, err := actor.NewFromKeyset(config.GetKeyset(), config.GetPublish())
	if err != nil {
		log.Warnf("Error initializing actor: %v", err)
	}

	fmt.Printf("I am : %s\n", a.Entity.DID.String())
	fmt.Printf("My public key is: %s\n", n.ID().String())

	// Now we can start continuous discovery in the background.
	go p.DiscoveryLoop(ctx)
	go handleEvents(ctx, a)

	// This is defined in web.go. It makes it possible to add extra parameters to the handler.
	h := &WebHandlerData{n, a}
	http.HandleFunc("/", h.WebHandler)
	log.Infof("Listening on %s\n", config.GetHttpSocket())
	err = http.ListenAndServe(config.GetHttpSocket(), nil)
	if err != nil {
		log.Fatal(err)
	}
}
