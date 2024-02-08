package ui

import (
	"context"

	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma/msg"
	p2ppubsub "github.com/libp2p/go-libp2p-pubsub"
	log "github.com/sirupsen/logrus"
)

// Handle incoming envelopes to an entity. The actor
// is responsible for decrypting the envelope. The entity
// Is only provided in order to decide whether to accept the message or not.
func (ui *ChatUI) handleIncomingEnvelopes(a *entity.Entity) {
	log.Debugf("Waiting for actor envelopes")
	for {
		envelope, ok := <-a.Envelopes // Envelopes should always have been sent to the actor.
		if !ok {
			log.Debug("Actor envelope channel closed, exiting...")
			return
		}
		log.Debugf("Received actor envelope: %v", envelope)

		// Process the envelope and send a pong response
		m, err := envelope.Open(a.Keyset.EncryptionKey.PrivKey[:])
		if err != nil {
			log.Errorf("Error opening actor envelope: %v\n", err)
			continue
		}

		log.Debugf("Opened envelope and found message: %v\n", string(m.Content))

		// Send the message to the actor for processing. It can decide to ignore it.
		a.Messages <- m
	}
}

// Subscribe to envelopes from an entity. But this time we are the entity
// to receive the envelopes. The actor is responsible for decrypting the envelope.
func (ui *ChatUI) subscribeToPubsubEnvelopes(e *entity.Entity, a *entity.Entity) {

	t := e.DID.String()

	ctx := context.Background()

	actor := a.DID.String()
	entity := e.DID.String()

	log.Debugf("Subscribing to envelopes to %s delivered to entity %s", actor, entity)
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
				log.Errorf("Actor message channel %s closed, exiting.... That's bad!", t)
				return
			}

			// If the message is not verified it might be an envelope.
			env, err := msg.UnmarshalAndVerifyEnvelopeFromCBOR(message.Data)
			if err != nil {
				log.Errorf("handleSubscriptionMessages: Failed to unmarshal and verify envelope: %v", err)
				continue
			}
			log.Debugf("handleSubscriptionMessages: Envelope verified: %v. Passing it on to actor %s", env, t)
			a.Envelopes <- env
		}
	}
}
