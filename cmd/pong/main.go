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
	initConfig(pong)

	// THese are the relay specific parts.

	p, err := p2p.Init(p2pOptions())
	if err != nil {
		fmt.Printf("Failed to initialize p2p: %v\n", err)
		return
	}

	// Init of actor requires P2P to be initialized
	a := actor.Init()

	fmt.Printf("Starting pong mode as %s\n", a.Entity.DID.Id)
	go p.StartDiscoveryLoop(ctx)
	fmt.Println("Discovery loop started.")
	go a.Subscribe(ctx, a.Entity)
	fmt.Println("Subscribed to self.")

	go handleEnvelopeEvents(ctx, a)
	go handleMessageEvents(ctx, a)
	fmt.Println("Started event handlers.")

	actor.HelloWorld(ctx, a)
	fmt.Println("Sent hello world.")

	fmt.Printf("Running in pong mode as %s@%s\n", a.Entity.DID.Id, p.Host.ID())
	fmt.Println("Press Ctrl-C to stop.")

	for {
		// A background loop that does nothing.
		// The ctx will never be cancelled, so this will run forever.
		<-ctx.Done()
		fmt.Println("Pong run loop cancelled, exiting...")
		return
	}
}
