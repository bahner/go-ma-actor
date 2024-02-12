package ui

import "context"

func (ui *ChatUI) startActor() {

	ctx := context.Background()

	// Now that the UI is created, we can start the actor and subscribe to its events.
	go ui.a.Subscribe(ctx, ui.a.Entity)

	// We *don't* want to subscribe to messages for the actor. We dont want to handle
	// messages for the actor. We want to handle envelopes for the actor.
	// Anyone can encrypt messages for the actor, so .. do that already!
	go ui.handleIncomingEnvelopes(ctx, ui.a)
	go ui.handleIncomingMessages(ctx, ui.a.Entity)

}
