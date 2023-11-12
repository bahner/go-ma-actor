package main

import (
	"context"
	"encoding/json"

	"github.com/bahner/go-ma/message"
	log "github.com/sirupsen/logrus"
)

func (a *Actor) Listen(ctx context.Context) {
	for {
		msg, err := a.From.Next(ctx)
		if err != nil {
			// Log the error or handle it more gracefully.
			log.Errorf("Failed to get next message: %v", err)
			return
		}

		if msg.ReceivedFrom == ps.Host.Node.ID() {
			continue
		}

		am := new(message.Message)
		if err := json.Unmarshal(msg.Data, am); err != nil {
			log.Debugf("Failed to unmarshal message: %v", err)
			continue
		}

		a.ProcessMessage(am)
	}
}

func (a *Actor) ProcessMessage(m *message.Message) {
	// Handle the message according to your application's logic.
	// For instance, this could involve updating the actor's state, triggering some action, etc.
	log.Debugf("process_message: Processed message: %s", m.ID)
}
