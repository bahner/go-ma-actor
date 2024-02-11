package ui

import (
	"context"

	"github.com/bahner/go-ma"
	"github.com/bahner/go-ma-actor/entity"
	log "github.com/sirupsen/logrus"
)

// Handle incoming messages to an entity. The topic is just to filter,
// which recipients *not* to handle here.
func (ui *ChatUI) handleIncomingMessages(ctx context.Context, e *entity.Entity) {
	t := e.Topic.String()
	log.Debugf("Waiting for messages from topic %s", t)

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

			// Only broadcasts to the actual subscriber.
			if m.MimeType == ma.BROADCAST_MIME_TYPE && m.To == t {
				log.Debugf("Received broadcast from %s to %s", m.From, m.To)
				ui.chMessage <- m
				continue
			}

			// Ignore messages to other topics than this goroutine's.
			if m.To != t {
				log.Debugf("Received message to %s. Expected %s. Ignoring...", m.To, t)
				continue
			}

			if m.From == t {
				log.Debugf("Received message from self, ignoring...")
				continue
			}

			log.Debugf("handleIncomingMessages: Accepted message of type %s from %s to %s", m.MimeType, m.From, m.To)
			ui.chMessage <- m
		}
	}
}
