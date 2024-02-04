package ui

import (
	log "github.com/sirupsen/logrus"
)

// Handle incoming messages to an entity
func (ui *ChatUI) handleIncomingMessages() {

	t := ui.e.Topic.String()

	log.Debugf("Waiting for messages from topic %s", t)

	for {
		m, ok := <-ui.e.Messages
		if !ok {
			log.Debug("Message channel closed, exiting...")
			return
		}
		log.Debugf("Received message from %s to %s", m.From, m.To)

		// Ignore self
		if m.From == ui.a.DID.String() {
			log.Debugf("Received message from self, ignoring...")
			continue
		}

		// Broadcast only allowed from the entity, ie. The Room.
		if m.From == m.To && m.From != t {
			log.Debugf("Received broadcast from %s, ignoring...", m.From)
			continue
		}

		log.Debugf("Received message from %s", m.From)
		ui.displayChatMessage(m)
	}
}
