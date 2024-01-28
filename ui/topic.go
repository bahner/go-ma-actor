package ui

import (
	"context"

	"github.com/bahner/go-ma-actor/p2p/topic"
	log "github.com/sirupsen/logrus"
)

// Handle incoming messages. NB! This must be cancelled,
// when topic (localtion/home) changes.
func (ui *ChatUI) handleTopicEvents() {

	envelopes := ui.t.SubscribeEnvelopes(ui.currentCtx)

	for {
		log.Debugf("Waiting for messages from topic %s", ui.t.Topic.String())
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

		case <-ui.currentCtx.Done():
			log.Debug("Context done, closing envelope channel...")
			return
		}
	}
}

func (ui *ChatUI) changeTopic(topicName string) {

	var err error

	// If there is an ongoing topic, cancel its context to stop the goroutine
	if ui.currentCancel != nil {
		log.Debugf("Cancelling current context")
		ui.currentCancel()
	}

	// Create a new context for the new topic
	ui.currentCtx, ui.currentCancel = context.WithCancel(context.Background())

	log.Debugf("Creating entity for topic %s", topicName)
	ui.e, err = getOrCreateEntity(topicName)
	if err != nil {
		log.Errorf("Failed to get or create entity: %v", err)
		return
	}

	// The channel for incoming messages
	log.Debugf("Creating topic for entity %s", ui.e.DID)
	ui.t, err = topic.GetOrCreate(ui.e.DID)
	if err != nil {
		log.Errorf("topic creation error: %s", err)
	}

	log.Infof("Topic changed to %s", ui.t.Topic.String())

	// Start handling the new topic
	go ui.handleTopicEvents()

}
