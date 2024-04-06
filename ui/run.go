package ui

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/entity"
)

// Run starts the chat event loop in the background, then starts
// the event loop for the text UI.
func (ui *ChatUI) Run() error {

	defer ui.end()

	// Now we can start continuous discovery in the background.
	fmt.Println("Starting discovery loop in the background....")
	go ui.p.StartDiscoveryLoop(context.Background())

	fmt.Println("Starting pubsub peers loop in the background....")
	go ui.pubsubPeersLoop(context.Background())

	// The actor should just run in the background for ever.
	// It will handle incoming messages and envelopes.
	// It shouldn't change - ever.
	fmt.Println("Starting actor...")
	ui.startActor()

	// We must wait for this to finish.
	fmt.Printf("Entering %s ...\n", config.ActorLocation())
	e, err := entity.GetOrCreate(config.ActorLocation())
	if err != nil {
		return fmt.Errorf("ui/run: %w", err)
	}
	err = ui.enterEntity(e, true)
	if err != nil {
		ui.displayStatusMessage(err.Error())
	}
	fmt.Printf("Entered %s\n", config.ActorLocation())

	// Subscribe to incoming messages
	go ui.handleIncomingMessages(context.Background())

	fmt.Println("Starting event loop...")
	go ui.handleEvents()

	return ui.app.Run()

}
