package ui

import (
	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma/msg"
	"github.com/fxamacker/cbor/v2"
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

// Subscribe a to e's topic and handle messages
// The envelopes are decrypted by ui.a - the actor. Not the entity.
// This must be called after the new entity is set.
func (ui *ChatUI) subscribeToEntityPubsubEnvelopes(a *entity.Entity) {

	t := a.DID.String()

	sub, err := a.Subscribe()
	if err != nil {
		log.Errorf("Failed to subscribe to topic: %v", err)
		return
	}

	defer sub.Cancel()

	for {
		log.Debugf("Waiting for pubsub messages for %s", t)
		input, ok := <-sub.Messages
		if !ok {
			log.Debugf("handleSubscriptionMessages: Input channel closed, exiting...")
			return
		}

		// Attempt to unmarshal the data into a msg.Envelope struct.
		var env *msg.Envelope
		err := cbor.Unmarshal(input.Data, &env)
		if err != nil {
			// If unmarshalling fails, log the error.
			log.Errorf("handleSubscriptionMessages: Error unmarshalling envelope: %v\n", err)
			// Here, you might want to return or continue based on your application's logic.
			// If this is not a critical error, you might choose to continue to try other data formats or handling.
			continue
		}

		// If unmarshalling succeeds, proceed to send the envelope to the actor.
		log.Debugf("handleSubscriptionMessages: Sending unmarshalled envelope to actor %s", a.DID.String())
		a.Envelopes <- env
	}
}
