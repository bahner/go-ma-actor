package actor

import (
	"context"

	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma/msg"
	p2ppubsub "github.com/libp2p/go-libp2p-pubsub"
	log "github.com/sirupsen/logrus"
)

const PUBSUB_MESSAGES_BUFFERSIZE = 32

// Subscribe to the entity's topic and handle incoming messages.
// the actor is the entity that will receive the messages.
// It may well be the entity itself.
// The context sho
func (a *Actor) Subscribe(ctx context.Context, e *entity.Entity) {

	// WHen an actor subscribes to an entity, it will receive messages and envelopes.
	// Messages should sent to the entity, whereas envelopes should be sent to the actor.

	they := e.DID.Id
	me := a.Entity.DID.Id

	log.Infof("Subscribing to %s as %s: ", they, me)

	// ctx, cancel := context.WithCancel(ctx)
	// defer cancel()

	sub, err := e.Topic.Subscribe()
	if err != nil {
		log.Errorf("Failed to subscribe to topic: %v", err)
		return
	}

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
			log.Debugf("Entity %s is cancelled, exiting subscription loop...", they)
			return
		case message, ok := <-messages:
			if !ok {
				log.Debugf("Message channel %s closed, exiting...", they)
				return
			}

			// Firstly check if this is a public message. Its quicker.
			m, err := msg.UnmarshalAndVerifyMessageFromCBOR(message.Data)
			if err == nil {
				log.Debugf("handleSubscriptionMessages: Received message: %v\n", m)
				e.Messages <- m
				continue
			} else {
				log.Debugf("handleSubscriptionMessages: Received message that is not a verified message: %v\n", err)
			}

			// If it's not a public message, it might be an envelope.
			env, err := msg.UnmarshalAndVerifyEnvelopeFromCBOR(message.Data)
			if err == nil {
				log.Debugf("handleSubscriptionMessages: Received envelope: %v\n", env)
				a.Envelopes <- env
				continue
			} else {
				log.Debugf("handleSubscriptionMessages: Received message that is not a verified envelope: %v\n", err)
			}

			log.Errorf("handleSubscriptionMessages: Received message that is neither a message nor an envelope: %v\n", message.Data)
		}
	}
}
