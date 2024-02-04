package ui

import (
	"github.com/bahner/go-ma-actor/entity"
	log "github.com/sirupsen/logrus"
)

// Handle incoming messages. NB! This must be cancelled,
// when topic (localtion/home) changes.
func (ui *ChatUI) handleIncomingEnvelopes(a *entity.Entity) {

	for {
		log.Debugf("Waiting for messages from topic %s", a.Topic.String())
		select {
		case e, ok := <-ui.e.Envelopes:
			if !ok {
				log.Debug("Envelope channel closed, exiting...")
				return
			}
			log.Debugf("Received envelope: %v", e)

			// Process the envelope and send a pong response
			m, err := e.Open(ui.a.Keyset.EncryptionKey.PrivKey[:])
			if err != nil {
				log.Errorf("Error opening envelope: %v\n", err)
				continue
			}

			log.Debugf("Received message: %v\n", string(m.Content))

			// Check if the message is from self to prevent pong loop
			if m.From != ui.a.DID.String() {
				log.Debugf("Received message from %s", m.From)
				ui.displayChatMessage(m)
			} else {
				log.Debugf("Received message from self, ignoring...")
			}

		case <-ui.currentCtx.Done():
			log.Debug("ui/handleIncomingEnvelopes, ui context done. Closing envelope channel...")
		}
	}
}
