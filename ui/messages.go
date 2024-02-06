package ui

import (
	"github.com/bahner/go-ma-actor/entity"
	log "github.com/sirupsen/logrus"
)

// Handle incoming messages to an entity. The topic is just to filter,
// which recipients *not* to handle here.
func (ui *ChatUI) handleIncomingMessages(e *entity.Entity) {

	t := e.Topic.String()

	log.Debugf("Waiting for messages from topic %s", t)

	for {
		m, ok := <-e.Messages
		if !ok {
			log.Debug("Message channel closed, exiting...")
			return
		}
		log.Debugf("Received message from %s to %s", m.From, m.To)

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
