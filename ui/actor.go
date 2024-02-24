package ui

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma-actor/internal"
	log "github.com/sirupsen/logrus"
)

func (ui *ChatUI) startActor() {

	log.Info("Starting actor...")
	ctx := context.Background()

	// Now that the UI is created, we can start the actor and subscribe to its events.
	fmt.Print("Subscribing to actor topic messages...")
	go ui.a.Subscribe(ctx, ui.a.Entity)
	fmt.Println("done.")

	// We want to handle envelopes for the actor, then deliver the messages
	// to the UI from the incoming envelopes.
	fmt.Print("Starting actor handleIncomingEnvelopes...")
	go ui.handleIncomingEnvelopes(ctx, ui.a.Entity, ui.a)
	fmt.Println("done.")
	fmt.Print("Starting actor handleIncomingMessages...")
	go ui.handleIncomingMessages(ctx, ui.a.Entity)
	fmt.Println("done.")

	internal.HelloWorld(ctx, ui.a, ui.b)
}
