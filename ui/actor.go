package ui

import "context"

func (ui *ChatUI) startActor() {

	ctx := context.Background()

	// Now that the UI is created, we can start the actor and subscribe to its events.
	go ui.a.Subscribe(ctx, ui.a.Entity)

	// We *don't* want to subscribe to messages for the actor.
	// We want to handle envelopes for the actor, then deliver the messages
	// to the UI from the incoming envelopes.
	go handleIncomingEnvelopes(ctx, ui.a.Entity, ui.a)
	go ui.handleIncomingMessages(ctx, ui.a.Entity)

}
