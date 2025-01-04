package main

import (
	"context"
	"fmt"
	"time"

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

	relayResources := relay.Resources{
		Limit: &relay.RelayLimit{
			Duration: 10 * time.Minute, // Allow longer connection durations
			Data:     10 << 20,         // Allow larger data transfer (10 MB)
		},
		MaxReservations: 256, // Allow more reservations
		MaxCircuits:     32,  // Allow more relayed streams per peer
	}

	opts := p2p.DefaultOptions()
	opts.P2P = append(opts.P2P,
		libp2p.EnableRelayService(),
		libp2p.ForceReachabilityPublic(),
		libp2p.EnableRelayService(relay.WithResources(relayResources)),
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
