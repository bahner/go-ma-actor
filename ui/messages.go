package ui

import (
	"context"

	"github.com/bahner/go-ma-actor/entity"
	log "github.com/sirupsen/logrus"
)

// Handle incoming messages to an entity. Also adds a reject to be able to filter self.
// If reject is nil, no filtering is done.
func (ui *ChatUI) handleIncomingMessages(ctx context.Context, e *entity.Entity) {
	t := e.Topic.String()
	me := ui.a.Entity.DID.Id

	log.Debugf("Handling incoming messages to %s", t)

	for {
		select {
		case <-ctx.Done():
			log.Debug("Context cancelled, exiting handleIncomingMessages...")
			return
		case m, ok := <-e.Messages:
			if !ok {
				log.Debug("Message channel closed, exiting...")
				return
			}
			log.Debugf("Received message from %s to %s", m.From, m.To)

			// Accept messages to the general topic or to the actor.
			if m.To == t || m.To == me {
				log.Debugf("handleIncomingMessages: Accepted message of type %s from %s to %s", m.MimeType, m.From, m.To)
				ui.chMessage <- m
				continue
			}

			log.Debugf("handleIncomingMessages: Received message to %s. Expected %s. Ignoring...", m.To, t)
		}
	}
}
