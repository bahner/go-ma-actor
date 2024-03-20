package ui

import (
	"context"

	log "github.com/sirupsen/logrus"
)

// Handle incoming messages to an entity. The message are recieved in the entity's message channel.
// And delivered to a channel of your choice.
func (ui *ChatUI) handleIncomingMessages(ctx context.Context) {

	log.Info("Handling incoming messages to UI...")

	for {
		select {
		case <-ctx.Done():
			log.Debug("Context cancelled, exiting ui.handleIncomingMessages...")
			return
		case m, ok := <-ui.chMessages:
			if !ok {
				log.Debug("Message channel closed, exiting...")
				return
			}
			content := string(m.Message.Content)
			from := m.Message.From
			to := m.Message.To
			log.Debugf("UI received message %v from %s to %s", content, from, to)

			// No need to verify at this point, as the message has already been verified by the actor.
			if m.Enveloped {
				ui.displayPrivateMessage(m.Message)
				continue
			}
			ui.displayChatMessage(m.Message)
		}
	}
}
