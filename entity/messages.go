package entity

import (
	"context"

	log "github.com/sirupsen/logrus"
)

// Handle incoming messages to an entity. The message are recieved in the e|ntity's message channel.
// And delivered to a channel of your choice.
func (e *Entity) HandleIncomingMessages(ctx context.Context, msgChan chan *Message) {
	me := e.DID.Id

	log.Debugf("Handling incoming messages to %s", me)

	for {
		select {
		case <-ctx.Done():
			log.Debug("Context cancelled, exiting entity.HandleIncomingMessages...")
			return
		case m, ok := <-e.Messages:
			if !ok {
				log.Debug("Message channel closed, exiting...")
				return
			}

			from := m.Message.From
			to := m.Message.To
			t := m.Message.Type

			log.Debugf("Received message from %s to %s", from, to)

			err := m.Message.Verify()
			if err != nil {
				log.Errorf("handleIncomingMessage: %s: %v", me, err)
				continue
			}

			if to == me {
				log.Debugf("entity.HandleIncomingMessages: Accepted message of type %s from %s to %s", t, from, to)
				msgChan <- m
				continue
			}

			log.Debugf("entity.HandleIncomingMessages: Received message to %s. Expected %s. Ignoring...", to, me)
		}
	}
}
