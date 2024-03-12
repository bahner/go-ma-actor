package entity

import (
	"context"

	"github.com/bahner/go-ma/msg"
	log "github.com/sirupsen/logrus"
)

// Handle incoming messages to an entity. The message are recieved in the entity's message channel.
// And delivered to a channel of your choice.
func (e *Entity) HandleIncomingMessages(ctx context.Context, msgChan chan<- *msg.Message) {
	me := e.DID.Id

	log.Debugf("Handling incoming messages to %s", me)

	for {
		select {
		case <-ctx.Done():
			log.Debug("Context cancelled, exiting handleIncomingMessages...")
			return
		case m, ok := <-e.Messages:
			if !ok {
				log.Debug("Message channel closed, exiting...")
				return
			}
			log.Debugf("Received message from %s to %s", m.From, m.To)

			err := m.Verify()
			if err != nil {
				log.Errorf("handleIncominngMessage: %s: %v", me, err)
				continue
			}

			if m.To == me {
				log.Debugf("handleIncomingMessages: Accepted message of type %s from %s to %s", m.Type, m.From, m.To)
				msgChan <- m
				continue
			}

			log.Debugf("handleIncomingMessages: Received message to %s. Expected %s. Ignoring...", m.To, me)
		}
	}
}
