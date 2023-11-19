package main

import (
	"context"
	"encoding/json"

	"github.com/bahner/go-ma/msg"
	"github.com/bahner/go-space/actor"
	log "github.com/sirupsen/logrus"
)

func (a *actor.Actor) Listen(ctx context.Context) {
	for {
		m, err := a.From.Next(ctx)
		if err != nil {
			// Log the error or handle it more gracefully.
			log.Errorf("Failed to get next message: %v", err)
			return
		}

		if m.ReceivedFrom == ps.Host.Node.ID() {
			continue
		}

		am := new(msg.Message)
		if err := json.Unmarshal(m.Data, am); err != nil {
			log.Debugf("Failed to unmarshal message: %v", err)
			continue
		}

		a.ProcessMessage(am)
	}
}

func (a *Actor) ProcessMessage(m *msg.Message) {
	// Handle the message according to your application's logic.
	// For instance, this could involve updating the actor's state, triggering some action, etc.
	log.Debugf("process_message: Processed message: %s", m.ID)
}
