package ui

import (
	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma/msg"
	"github.com/fxamacker/cbor/v2"
	log "github.com/sirupsen/logrus"
)

// Subscribe a to e's topic and handle messages
// The envelopes are decrypted by ui.a - the actor. Not the entity.
// This must be called after the new entity is set.
func (ui *ChatUI) subscribeToEntityPubsubMessages(e *entity.Entity) {

	t := e.DID.String()

	// log.Debugf("Subscribing to entity %s", e.DID.String())
	// sub, err := e.Subscribe()
	// if err != nil {
	// 	log.Errorf("Failed to subscribe to topic: %v", err)
	// 	return
	// }
	// defer sub.Cancel()

	for {
		log.Debugf("Waiting for pubsub messages to entity %s", t)
		input, ok := <-e.Subscription.Messages
		if !ok {
			log.Debugf("handleSubscriptionMessages: Input channel closed, exiting...")
			return
		}

		var m *msg.Message
		err := cbor.Unmarshal(input.Data, &m)
		if err != nil {
			// If unmarshalling fails, log the error and possibly continue or return.
			log.Errorf("handleSubscriptionMessages: Error unmarshalling message: %v\n", err)
			continue
		}

		// Log the received message.
		log.Debugf("handleSubscriptionMessages: Received message: %v\n", m)

		// Discard message if this isn't the correct topic.
		if m.To != t {
			log.Debugf("handleSubscriptionMessages: Received message to %s. Expected %s. Ignoring...", m.To, t)
			continue
		}

		// Verify the message.
		err = m.Verify()
		if err != nil {
			log.Debugf("handleSubscriptionMessages: Message verification failed: %v\n", err)
			continue
		}

		log.Debugf("handleSubscriptionMessages: Message verified: %v\n", m)
		e.Messages <- m
	}
}
