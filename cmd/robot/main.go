package main

import (
	"context"
	"fmt"
	"log"

	"github.com/bahner/go-ma-actor/ui/web"

	"github.com/bahner/go-ma-actor/p2p"
)

const defaultProfileName = "robot"

// Run the pong actor. Cancel it from outside to stop it.
func main() {

	ctx := context.Background()
	initConfig(defaultProfileName)

	i, err := NewRobot()
	if err != nil {
		log.Fatal(err)
	}

	p, err := p2p.Init(p2p.DefaultP2POptions())
	if err != nil {
		fmt.Printf("Failed to initialize p2p: %v\n", err)
		return
	}

	go p.StartDiscoveryLoop(ctx)

	i.Robot.HelloWorld(ctx)
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
