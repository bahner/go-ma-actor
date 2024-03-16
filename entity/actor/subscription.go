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
// Takes a channel for message delivery.
func (a *Actor) Subscribe(ctx context.Context, e *entity.Entity) {

	// WHen an actor subscribes to an entity, it will receive messages and envelopes.
	// Messages should sent to the entity, whereas envelopes should be sent to the actor.

	them := e.DID.Id
	me := a.Entity.DID.Id

	log.Infof("Subscribing to %s as %s: ", them, me)

	// ctx, cancel := context.WithCancel(ctx)
	// defer cancel()

	sub, err := e.Topic.Subscribe()
	if err != nil {
		log.Errorf("Failed to subscribe to topic: %v", err)
		return
	}
	cancelRelay, err := e.Topic.Relay()
	if err != nil {
		log.Errorf("actorSubscribe: failed to relay to topic: %v", err)
	}
	defer cancelRelay()

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
			log.Debugf("actor.Subscribe: Received message: %s", message.Data)

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
			log.Debugf("actor.Subscribe: Entity %s is cancelled, exiting subscription loop...", them)
			return
		case message, ok := <-messages:
			if !ok {
				log.Debugf("actor.Subscribe: Message channel %s closed, exiting...", them)
				return
			}

			// Firstly check if this is a public message. Its quicker.
			m, err := msg.UnmarshalAndVerifyMessageFromCBOR(message.Data)
			if err == nil && m != nil {
				log.Debugf("actor.Subscribe: Received message %s to %s\n", m.Id, m.To)
				log.Debugf("actor.Subscribe: Delivering message %s to entity: %s\n", m.Id, them)
				if m.To == me {
					log.Debugf("actor.Subscribe: Received message for actor: %s\n", me)
					a.Entity.Messages <- m
					continue
				}

				if m.To == them {
					e.Messages <- m
					continue
				}

				log.Debugf("actor.Subscribe: Received message to %s. Expected %s or %s. Ignoring...", m.To, me, them)
			} else {
				log.Debugf("actor.Subscribe: Received message that is not a verified message: %v\n", err)
			}

			// If it's not a public message, it might be an envelope.
			env, err := msg.UnmarshalAndVerifyEnvelopeFromCBOR(message.Data)
			if err == nil {
				log.Debugf("actor.Subscribe: Received envelope: %v\n", env)
				log.Debugf("actor.Subscribe: Delivering envelope to actor: %s\n", me)
				a.Envelopes <- env
				continue
			} else {
				log.Debugf("actor.Subscribe: Received message that is not a verified envelope: %v\n", err)
			}

			log.Errorf("actor.Subscribe: Received message that is neither a message nor an envelope: %v\n", message.Data)
		}
	}
}
