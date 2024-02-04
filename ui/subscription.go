package ui

import (
	"encoding/json"

	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma/msg"
	"github.com/fxamacker/cbor/v2"
	log "github.com/sirupsen/logrus"
)

// Subscribe a to e's topic and handle messages
// The envelopes are decrypted by ui.a - the actor. Not the entity.
func (ui *ChatUI) subscribeEntityMessages(e *entity.Entity) {
	sub, err := e.Subscribe()
	if err != nil {
		log.Errorf("Failed to subscribe to topic: %v", err)
		return
	}
	defer sub.Cancel()

	go ui.handleIncomingMessages(ui.a.DID.String())

	for {
		input, ok := <-sub.Messages
		if !ok {
			log.Debugf("handleSubscriptionMessages: Input channel closed, exiting...")
			return
		}

		// Firstly check if this is a public message. Its quicker.
		// Sent them directly to the entity.
		// Attempt to unmarshal the data into a msg.Message struct.
		var m *msg.Message
		err := cbor.Unmarshal(input.Data, &m)
		if err != nil {
			// If unmarshalling fails, log the error and possibly continue or return.
			log.Errorf("handleSubscriptionMessages: Error unmarshalling message: %v\n", err)
			continue
		}

		// Log the received message.
		log.Debugf("handleSubscriptionMessages: Received message: %v\n", m)

		// Verify the message.
		if err := m.Verify(); err == nil {
			log.Debugf("handleSubscriptionMessages: Message verified: %v\n", m)
			ui.e.Messages <- m
			// Successfully verified message, so continue to the next iteration.
			continue
		} else {
			// If verification fails, log the verification error.
			log.Debugf("handleSubscriptionMessages: Message verification failed: %v\n", err)
			// Optionally marshal the message to JSON for logging, if needed.
			msgJson, jsonErr := json.Marshal(m)
			if jsonErr != nil {
				log.Debugf("handleSubscriptionMessages: Error marshalling message to JSON: %v\n", jsonErr)
			} else {
				log.Debugf("handleSubscriptionMessages: Message not verified: %s\n", string(msgJson))
			}
		}

		// Envelopes goes to the actor, not the entity
		// Attempt to unmarshal the data into a msg.Envelope struct.
		var env *msg.Envelope
		err = cbor.Unmarshal(input.Data, &env)
		if err != nil {
			// If unmarshalling fails, log the error.
			log.Errorf("handleSubscriptionMessages: Error unmarshalling envelope: %v\n", err)
			// Here, you might want to return or continue based on your application's logic.
			// If this is not a critical error, you might choose to continue to try other data formats or handling.
			continue
		}

		// If unmarshalling succeeds, proceed to send the envelope to the actor.
		log.Debugf("handleSubscriptionMessages: Sending envelope to actor %s", ui.a.DID.String())
		ui.a.Envelopes <- env
	}
}
