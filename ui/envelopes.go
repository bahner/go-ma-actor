package ui

import (
	log "github.com/sirupsen/logrus"
)

// Handle incoming envelopes to an entity. The actor
// is responsible for decrypting the envelope. The entity
// Is only provided in order to decide whether to accept the message or not.
func (ui *ChatUI) handleIncomingEnvelopes() {

	log.Debugf("Waiting for actor envelopes")
	for {
		e, ok := <-ui.a.Envelopes // Envelopes should always have been sent to the actor.
		if !ok {
			log.Debug("Actor envelope channel closed, exiting...")
			return
		}
		log.Debugf("Received actor envelope: %v", e)

		// Process the envelope and send a pong response
		m, err := e.Open(ui.a.Keyset.EncryptionKey.PrivKey[:])
		if err != nil {
			log.Errorf("Error opening actor envelope: %v\n", err)
			continue
		}

		log.Debugf("Open actor envelope and found message: %v\n", string(m.Content))

		// Check if the message is from self to prevent loop
		if m.From == ui.a.DID.String() {
			log.Debugf("Received message from self(%s), ignoring...", m.From)
			continue
		}

		log.Debugf("Displaying message from %s", m.From)
		ui.displayChatMessage(m)
	}
}
