package main

import (
	"context"
	"fmt"
	"log"

	"github.com/bahner/go-ma-actor/entity/actor"
	"github.com/bahner/go-ma-actor/ui/web"

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

	go p.StartDiscoveryLoop(ctx)

	i, err := NewRobot()
	if err != nil {
		log.Fatal(err)
	}

	actor.HelloWorld(ctx, i.Robot)
	// i.Robot.HelloWorld(ctx, a)

	fmt.Println("Press Ctrl-C to stop.")

	h := web.NewEntityHandler(p, i.Robot.Entity)
	go web.Start(h)

	for {
		// A background loop that does nothing.
		// The ctx will never be cancelled, so this will run forever.
		<-ctx.Done()
		fmt.Println("Pong run loop cancelled, exiting...")
		return
	}
}
