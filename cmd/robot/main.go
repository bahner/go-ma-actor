package main

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma-actor/entity/actor"

	"github.com/bahner/go-ma-actor/p2p"
)

// Run the pong actor. Cancel it from outside to stop it.
func main() {

	ctx := context.Background()
	initConfig(name)

	// THese are the relay specific parts.

	p, err := p2p.Init(p2p.DefaultOptions())
	if err != nil {
		fmt.Printf("Failed to initialize p2p: %v\n", err)
		return
	}

	// Init of actor requires P2P to be initialized
	a := actor.Init()

	go p.StartDiscoveryLoop(ctx)

	// Subscribe self
	go a.Subscribe(ctx, a.Entity)
	go handleEnvelopeEvents(ctx, a)
	go handleMessageEvents(ctx, a)

	actor.HelloWorld(ctx, a)

	fmt.Println("Press Ctrl-C to stop.")

	for {
		// A background loop that does nothing.
		// The ctx will never be cancelled, so this will run forever.
		<-ctx.Done()
		fmt.Println("Pong run loop cancelled, exiting...")
		return
	}
}
