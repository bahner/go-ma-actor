package ui

import (
	"context"

	"github.com/bahner/go-ma/msg"
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

	greeting := []byte("Hello, world! " + ui.a.Entity.DID.Fragment + " is here.")
	mesg, err := msg.NewBroadcast(ui.a.Entity.DID.Id, greeting, "text/plain", ui.a.Keyset.SigningKey.PrivKey)
	if err != nil {
		ui.displaySystemMessage("Error creating greeting message: " + err.Error())
	}

	mesg.Send(ctx, ui.a.Entity.Topic)

}
