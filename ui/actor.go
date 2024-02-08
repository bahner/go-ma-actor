package ui

import "context"

func (ui *ChatUI) startActor() {

	ctx := context.Background()

	// Now that the UI is created, we can start the actor event loop.
	// There is no need to subscribe to messages for the actor.
	// Envelopes are handled by the actor.
	go ui.subscribeToEntityPubsubMessages(ctx, ui.a)
	// Subscribe to envelopes for the actor directly.
	go ui.subscribeToPubsubEnvelopes(ui.a, ui.a)

	// We *don't* want to subscribe to messages for the actor. We dont want to handle
	// messages for the actor. We want to handle envelopes for the actor.
	// Anyone can encrypt messages for the actor, so .. do that already!
	go ui.handleIncomingEnvelopes(ui.a)
	go ui.handleIncomingMessages(ui.a)

}
