package ui

import (
	"github.com/bahner/go-ma-actor/entity"
	log "github.com/sirupsen/logrus"
)

// Handle incoming messages. NB! This must be cancelled,
// when topic (localtion/home) changes.
func (ui *ChatUI) handleIncomingMessages(a *entity.Entity) {

	for {
		log.Debugf("Waiting for messages from topic %s", a.Topic.String())
		select {
		case m, ok := <-ui.e.Messages:
			if !ok {
				log.Debug("Message channel closed, exiting...")
				return
			}
			log.Debugf("Received message from %s to %s", m.From, m.To)

			// Check if the message is sent to the topic.
			if m.To == a.Topic.String() {
				log.Debugf("Received message from %s", m.From)
				ui.displayChatMessage(m)
			} else {
				log.Debugf("Ignoring message to %s in %s", m.To, a.Topic.String())
			}

		case <-ui.currentCtx.Done():
			log.Debug("ui/handleIncomingMessage, ui context done. Closing envelope channel...")
			return
		}
	}
}
