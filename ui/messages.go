package ui

import (
	"context"

	log "github.com/sirupsen/logrus"
)

// Handle incoming messages to an entity. The message are recieved in the entity's message channel.
// And delivered to a channel of your choice.
func (ui *ChatUI) handleIncomingMessages(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Debug("Context cancelled, exiting handleIncomingMessages...")
			return
		case m, ok := <-ui.chMessages:
			if !ok {
				log.Debug("Message channel closed, exiting...")
				return
			}
			log.Debugf("UI received message from %s to %s", m.From, m.To)

			// No need to verify at this point, as the message has already been verified by the actor.
			ui.displayChatMessage(m)
		case m, ok := <-ui.chPrivateMessages:
			if !ok {
				log.Debug("Private message channel closed, exiting...")
				return
			}
			log.Debugf("UI received private message from %s to %s", m.From, m.To)

			// No need to verify at this point, as the message has already been verified by the actor.
			ui.displayPrivateMessage(m)
		}
	}
}
