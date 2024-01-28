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

	config.Init("home")
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

	n := p.Node

	a, err := actor.NewFromKeyset(config.GetKeyset(), config.GetPublish())
	if err != nil {
		log.Warnf("Error initializing actor: %v", err)
	}

	fmt.Printf("I am : %s\n", a.Entity.DID.String())
	fmt.Printf("My public key is: %s\n", n.ID().String())

	go p.DiscoverPeers()
	go handleEvents(ctx, a)

	// This is defined in web.go. It makes it possible to add extra parameters to the handler.
	h := &WebHandlerData{n, a}
	http.HandleFunc("/", h.WebHandler)
	log.Infof("Listening on %s\n", getHttpSocket())
	err = http.ListenAndServe(getHttpSocket(), nil)
	if err != nil {
		log.Fatal(err)
	}
}
