package main

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/db"
	"github.com/bahner/go-ma-actor/p2p"
	"github.com/bahner/go-ma-actor/ui/web"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
)

const defaultProfileName = "relay"

// Run the pong actor. Cancel it from outside to stop it.
func main() {

	ctx := context.Background()
	initConfig(defaultProfileName)

	identity, err := db.GetOrCreateIdentity(config.Profile())
	if err != nil {
		fmt.Printf("Failed to get or create identity: %v\n", err)
		panic(err)
	}

	opts := p2p.DefaultOptions()
	opts.P2P = append(opts.P2P,
		libp2p.EnableRelay(),        // Enable relay support for relayed connections
		libp2p.EnableRelayService(), // Allow acting as a relay server
	)

	p, err := p2p.Init(identity, opts)
	if err != nil {
		fmt.Printf("Failed to initialize p2p: %v\n", err)
		return
	}

	_, err = relay.New(p.Host)
	if err != nil {
		fmt.Printf("Failed to initialize relay: %v\n", err)
		return
	}

	go p.StartDiscoveryLoop(ctx)
	fmt.Println("Discovery loop started.")

	handler := web.NewRelayHandler(p)
	web.Start(handler)

}
