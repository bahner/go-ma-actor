package main

import (
	"context"
	"fmt"
	"time"

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

	relayResources := relay.Resources{
		Limit: &relay.RelayLimit{
			Duration: 10 * time.Minute, // Allow longer connection durations
			Data:     10 << 20,         // Allow larger data transfer (10 MB)
		},
		MaxReservations: 256, // Allow more reservations
		MaxCircuits:     32,  // Allow more relayed streams per peer
	}

	p2pOpts := p2p.DefaultP2POptions()
	p2pOpts.P2P = append(p2pOpts.P2P,
		libp2p.EnableRelayService(),
		libp2p.ForceReachabilityPublic(),
		libp2p.EnableRelayService(relay.WithResources(relayResources)),
	)

	p, err := p2p.Init(p2pOpts)
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
