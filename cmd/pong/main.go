package main

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma-actor/entity/actor"
	"github.com/bahner/go-ma-actor/ui/web"

	"github.com/bahner/go-ma-actor/p2p"
)

const (
	defaultPongReply   = "Pong!"
	defaultFortuneMode = false
	pong               = "pong"
	defaultProfileName = pong
	defaultFortuneArgs = "-s"
)

// Run the pong actor. Cancel it from outside to stop it.
func main() {

	ctx := context.Background()

	initConfig(defaultProfileName)

	// Init of actor requires P2P to be initialized
	a := actor.Init()

	// THese are the relay specific parts.
	p, err := p2p.Init(a.Keyset.Identity, p2p.DefaultOptions())
	if err != nil {
		fmt.Printf("Failed to initialize p2p: %v\n", err)
		return
	}

	fmt.Printf("Starting pong mode as %s\n", a.Entity.DID.Id)
	go p.StartDiscoveryLoop(ctx)
	fmt.Println("Discovery loop started.")
	go a.Subscribe(ctx, a.Entity)
	fmt.Println("Subscribed to self.")

	go handleEnvelopeEvents(ctx, a)
	go handleMessageEvents(ctx, a)
	fmt.Println("Started event handlers.")

	a.HelloWorld(ctx)
	fmt.Println("Sent hello world.")

	// WEB
	fmt.Println("Initialising web UI...")
	wh := web.NewEntityHandler(p, a.Entity)
	go web.Start(wh)

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
