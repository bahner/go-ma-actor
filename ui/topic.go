package ui

import (
	"context"

	"github.com/bahner/go-ma-actor/p2p/topic"
	log "github.com/sirupsen/logrus"
)

// Handle incoming messages. NB! This must be cancelled,
// when topic (localtion/home) changes.
func (ui *ChatUI) handleTopicEvents() {

	ctx := context.Background()

	t, err := topic.GetOrCreate(ui.e.DID)
	envelopes := t.SubscribeEnvelopes(ctx)
	if err != nil {
		log.Debugf("topic creation error: %s", err)
		return
	}

	for {
		log.Debugf("Waiting for messages...")
		select {
		case e, ok := <-envelopes:
			if !ok {
				log.Debug("Envelope channel closed, exiting...")
				return
			}
			log.Debugf("Received envelope: %v", e)

			// Process the envelope and send a pong response
			m, err := e.Open(ui.a.Entity.Keyset.EncryptionKey.PrivKey[:])
			if err != nil {
				log.Errorf("Error opening envelope: %v\n", err)
				continue
			}

			log.Debugf("Received message: %v\n", string(m.Content))

			// Check if the message is from self to prevent pong loop
			if m.From != ui.a.Entity.DID.String() {
				log.Debugf("Received message from %s", m.From)
				ui.displayChatMessage(m)
			} else {
				log.Debugf("Received message from self, ignoring...")
			}

		case <-ctx.Done():
			log.Debug("Context done, exiting...")
			return
		}
	}
}
