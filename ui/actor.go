package ui

func (ui *ChatUI) startActor() {

	// Now that the UI is created, we can start the actor event loop.
	// There is no need to subscribe to messages for the actor.
	// Envelopes are handled by the actor.
	go ui.subscribeToActorPubsubMessages(ui.a)
	go ui.subscribeToActorPubsubEnvelopes(ui.a)

	// We *don't* want to subscribe to messages for the actor. We dont want to handle
	// messages for the actor. We want to handle envelopes for the actor.
	// Anyone can encrypt messages for the actor, so .. do that already!
	go ui.handleIncomingEnvelopes(ui.a)

}
