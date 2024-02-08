package ui

import (
	"context"

	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma/msg"
	p2ppubsub "github.com/libp2p/go-libp2p-pubsub"
	log "github.com/sirupsen/logrus"
)

// Handle incoming messages to an entity. The topic is just to filter,
// which recipients *not* to handle here.
func (ui *ChatUI) handleIncomingMessages(e *entity.Entity) {

	t := e.Topic.String()

	log.Debugf("Waiting for messages from topic %s", t)

	for {
		m, ok := <-e.Messages
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

		if m.From == t {
			log.Debugf("Received message from self, ignoring...")
			continue
		}

		log.Debugf("handleIncomingMessages: Accepted message of type %s from %s to %s", m.MimeType, m.From, m.To)

		ui.chMessage <- m
	}
}

func (ui *ChatUI) subscribeToEntityPubsubMessages(e *entity.Entity) {

	t := e.DID.String()

	log.Debugf("Subscribing to entity %s", e.DID.String())
	sub, err := e.Topic.Subscribe()
	if err != nil {
		log.Errorf("Failed to subscribe to topic: %v", err)
		return
	}
	defer sub.Cancel()

	// Create a channel for messages.
	messages := make(chan *p2ppubsub.Message, PUBSUB_MESSAGES_BUFFERSIZE)

	// Start an anonymous goroutine to bridge sub.Next() to the messages channel.
	go func() {
		for {
			// Assuming sub.Next() blocks until the next message is available,
			// and returns a message or an error.
			message, err := sub.Next(ui.currentEntityCtx)
			if err != nil {
				// Handle error (e.g., log, break, or continue based on the error type).
				log.Errorf("Error getting next message: %v", err)
				return // or continue based on your error handling policy
			}
			log.Debugf("handleSubscriptionMessages: Received message: %s", message.Data)

			// Assuming message is of the type you need; otherwise, adapt as necessary.
			select {
			case messages <- message:
			case <-ui.currentEntityCtx.Done():
				return
			}
		}
	}()

	for {
		select {
		case <-ui.currentEntityCtx.Done():
			log.Debugf("Entity %s is cancelled, exiting subscription loop...", t)
			return
		case message, ok := <-messages:
			if !ok {
				log.Debugf("Message channel %s closed, exiting...", t)
				return
			}

			m, err := msg.UnmarshalAndVerifyMessageFromCBOR(message.Data)
			if err == nil {
				e.Messages <- m
				log.Debugf("handleSubscriptionMessages: Message verified: %v. Passing it on to %s", m, t)
				continue
			}

			log.Errorf("handleSubscriptionMessages: Failed to verify message or envelope: %v", err)
		}
	}
}

func (ui *ChatUI) subscribeToActorPubsubMessages(e *entity.Entity) {

	t := e.DID.String()

	ctx := context.Background()

	log.Debugf("Subscribing to entity %s", e.DID.String())
	sub, err := e.Topic.Subscribe()
	if err != nil {
		log.Errorf("Failed to subscribe to topic: %v", err)
		return
	}
	defer sub.Cancel()

	// Create a channel for messages.
	messages := make(chan *p2ppubsub.Message, PUBSUB_MESSAGES_BUFFERSIZE)

	// Start an anonymous goroutine to bridge sub.Next() to the messages channel.
	go func() {
		for {
			// Assuming sub.Next() blocks until the next message is available,
			// and returns a message or an error.
			message, err := sub.Next(ctx)
			if err != nil {
				// Handle error (e.g., log, break, or continue based on the error type).
				log.Errorf("Error getting next message: %v", err)
				return // or continue based on your error handling policy
			}
			log.Debugf("handleSubscriptionMessages: Received message: %s", message.Data)

			// Assuming message is of the type you need; otherwise, adapt as necessary.
			select {
			case messages <- message:
			case <-ctx.Done():
				return
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			log.Debugf("Entity %s is cancelled, exiting subscription loop...", t)
			return
		case message, ok := <-messages:
			if !ok {
				log.Debugf("Message channel %s closed, exiting...", t)
				return
			}

			m, err := msg.UnmarshalAndVerifyMessageFromCBOR(message.Data)
			if err == nil {
				e.Messages <- m
				log.Debugf("handleSubscriptionMessages: Message verified: %v. Passing it on to %s", m, t)
				continue
			}

			log.Errorf("handleSubscriptionMessages: Failed to verify message or envelope: %v", err)
		}
	}
}
