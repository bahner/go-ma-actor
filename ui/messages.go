package ui

import (
	"github.com/bahner/go-ma"
	log "github.com/sirupsen/logrus"
)

// Handle incoming messages to an entity. The topic is just to filter,
// which recipients *not* to handle here.
func (ui *ChatUI) handleIncomingMessages(t string) {

	log.Debugf("Waiting for messages from topic %s", t)

	for {
		m, ok := <-ui.e.Messages
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

		if m.From == ui.a.DID.String() {
			log.Debugf("Received message from self, ignoring...")
			continue
		}

		// Handle broadcast messages
		// Allow broadcast sent to the topic
		if m.MimeType == ma.BROADCAST_MIME_TYPE && m.To != t {
			log.Debugf("Received broadcast from %s to %s, ignoring...", m.From, t)
			continue
		}

		log.Debugf("Received message from %s", m.From)
		ui.displayChatMessage(m)
	}
}
