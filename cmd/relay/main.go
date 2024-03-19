package main

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma-actor/p2p"
	"github.com/bahner/go-ma-actor/ui/web"
)

const relay = "relay"

// Run the pong actor. Cancel it from outside to stop it.
func main() {

	ctx := context.Background()
	initConfig(relay)

	p, err := p2p.Init(p2pOptions())
	if err != nil {
		fmt.Printf("Failed to initialize p2p: %v\n", err)
		return
	}

	go p.StartDiscoveryLoop(ctx)
	fmt.Println("Discovery loop started.")

	handler := web.NewRelayHandler(p)
	web.Start(handler)

}
