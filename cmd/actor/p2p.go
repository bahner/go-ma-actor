package main

import (
	"fmt"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/mode/pong"
	"github.com/bahner/go-ma-actor/mode/relay"
	"github.com/bahner/go-ma-actor/p2p"
	"github.com/bahner/go-ma-actor/p2p/connmgr"
)

func initP2P() (P2P *p2p.P2P, err error) {
	fmt.Println("Initialising libp2p...")

	// Everyone needs a connection manager.
	cm, err := connmgr.Init()
	if err != nil {
		panic(fmt.Errorf("pong: failed to create connection manager: %w", err))
	}
	cg := connmgr.NewConnectionGater(cm)

	if config.RelayMode() {
		fmt.Println("Relay mode enabled.")
		d, err := relay.DHT(cg)
		if err != nil {
			panic(fmt.Sprintf("failed to initialize dht: %v", err))
		}
		return p2p.Init(d)
	}

	if config.PongMode() {
		fmt.Println("Pong mode enabled.")
		d, err := pong.DHT(cg)
		if err != nil {
			panic(fmt.Sprintf("failed to initialize dht: %v", err))
		}
		return p2p.Init(d)
	}

	fmt.Println("Actor mode enabled.")
	d, err := DHT(cg)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize dht: %v", err))
	}
	return p2p.Init(d)
}
