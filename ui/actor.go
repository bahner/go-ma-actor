package ui

import (
	"context"

	log "github.com/sirupsen/logrus"
)

func (ui *ChatUI) startActor() {

	log.Info("Starting actor...")
	ctx := context.Background()

	// Now that the UI is created, we can start the actor and subscribe to its events.
	ui.a.Subscribe(ctx, ui.a.Entity)
	ui.displayStatusMessage("Subscribed to actor events")

	// We *don't* want to subscribe to messages for the actor.
	// We want to handle envelopes for the actor, then deliver the messages
	// to the UI from the incoming envelopes.
	go ui.handleIncomingEnvelopes(ctx, ui.a.Entity, ui.a)
	go ui.handleIncomingMessages(ctx, ui.a.Entity)

}
