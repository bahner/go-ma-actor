package main

import (
	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma/msg"
	"github.com/fxamacker/cbor/v2"
	log "github.com/sirupsen/logrus"
)

// SUbscribe a to e's topic and handle messages
func handleSubscriptionMessages(a *entity.Entity) {
	sub, err := a.Subscribe()
	if err != nil {
		log.Errorf("Failed to subscribe to topic: %v", err)
		return
	}
	defer sub.Cancel()

	for {
		select {
		// Handle cancellation
		case input, ok := <-sub.Messages:
			if !ok {
				log.Debugf("handleSubscriptionMessages: Input channel closed, exiting...")
				return
			}

			// Firstly check if this is a public message. Its quicker.
			var m *msg.Message
			err := cbor.Unmarshal(input.Data, &m)
			if err == nil {
				log.Debugf("handleSubscriptionMessages: Received message: %v\n", m)
				a.Messages <- m
				continue
			}

			var env *msg.Envelope
			err = cbor.Unmarshal(input.Data, &env)
			if err == nil {
				a.Envelopes <- env
				continue
			}
			log.Errorf("handleSubscriptionMessages: Error unmarshalling envelope: %v\n", err)
		case envelope, ok := <-a.Envelopes:
			if !ok {
				log.Errorf("handleSubscriptionMessages: Envelope channel closed, exiting...")
				continue
			}

			if a.Keyset.EncryptionKey == nil {
				log.Errorf("handleSubscriptionMessages: No encryption key found, cannot open envelope")
				continue
			}
			msg, err := envelope.Open(a.Keyset.EncryptionKey.PrivKey[:])
			if err == nil {
				log.Debugf("handleSubscriptionMessages: Open envelope: %v\n", msg)
				a.Messages <- msg
			}
		}
	}
}
