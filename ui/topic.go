package ui

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma-actor/p2p/topic"
	log "github.com/sirupsen/logrus"
)

// Handle incoming messages. NB! This must be cancelled,
// when topic (localtion/home) changes.
func (ui *ChatUI) handleTopicEvents(ctx context.Context, t *topic.Topic) {

	envelopes := t.SubscribeEnvelopes(ctx)

	for {
		log.Debugf("Waiting for messages from topic %s", t.Topic.String())
		select {
		case e, ok := <-envelopes:
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
			log.Debug("Context done, closing envelope channel...")
			return
		}
	}
}

func (ui *ChatUI) changeEntity(did string) error {

	var err error

	log.Debugf("Creating entity for topic %s", did)
	// e, err = getOrCreateEntity(did)
	e, err := entity.GetOrCreate(did)
	if err != nil {
		return fmt.Errorf("error getting or creating entity: %w", err)
	}

	// Now pivot to the new entity
	old_entity := ui.e
	ui.e = e
	old_entity.Leave()

	log.Infof("Location changed to %s", ui.e.Topic.Topic.String())

	// Start handling the new topic
	go ui.handleTopicEvents(ui.currentCtx, ui.e.Topic)

	return nil

}
