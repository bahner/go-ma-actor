package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/bahner/go-ma-actor/actor"
	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/p2p"

	log "github.com/sirupsen/logrus"
)

// var (
// 	ctx context.Context
// 	err error

// 	a *actor.Actor
// 	p *p2p.P2P
// 	n host.Host

// 	envelopes <-chan *msg.Envelope
// )

func main() {

	flag.Parse()
	config.Init()

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
		log.Fatalf("Error initializing actor: %v", err)
	}

	fmt.Printf("I am : %s\n", a.Entity.DID.String())
	fmt.Printf("My public key is: %s\n", n.ID().String())

	go discoveryHandler(p)
	go handleEvents(ctx, a)

	// This is defined in web.go. It makes it possible to add extra parameters to the handler.
	h := &WebHandlerData{n, a}
	http.HandleFunc("/", h.WebHandler)
	fmt.Println("Listening on port 5003...")
	err = http.ListenAndServe("0.0.0.0:5003", nil)
	if err != nil {
		log.Fatal(err)
	}
}
