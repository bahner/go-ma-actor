package main

import (
	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma/msg"
	"github.com/fxamacker/cbor/v2"
	log "github.com/sirupsen/logrus"
)

func handleSubscriptionMessages(e *entity.Entity) {
	sub, err := e.Enter(e)
	if err != nil {
		log.Errorf("Failed to subscribe to topic: %v", err)
		return
	}

	for {
		select {
		// Handle cancellation
		case <-e.Ctx.Done():
			log.Errorf("pong/handleSubscriptionMessages: Entity context done, exiting...")
			return
		case input, ok := <-sub:
			if !ok {
				log.Debugf("pong/handleSubscriptionMessages: Input channel closed, exiting...")
				return
			}

			// Firstly check if this is a public message. Its quicker.
			var m *msg.Message
			err := cbor.Unmarshal(input.Data, &m)
			if err == nil {
				log.Debugf("pong/handleSubscriptionMessages:Received message: %v\n", m)
				e.Messages <- m
				continue
			}

			var env *msg.Envelope
			err = cbor.Unmarshal(input.Data, &env)
			if err == nil {
				e.Envelopes <- env
				continue
			}
			log.Errorf("pong/handleSubscriptionMessages: Error unmarshalling envelope: %v\n", err)
		case envelope, ok := <-e.Envelopes:
			if !ok {
				log.Errorf("pong/handleSubscriptionMessages: Envelope channel closed, exiting...")
				continue
			}

			if e.Keyset.EncryptionKey == nil {
				log.Errorf("pong/handleSubscriptionMessages: No encryption key found, cannot open envelope")
				continue
			}
			msg, err := envelope.Open(e.Keyset.EncryptionKey.PrivKey[:])
			if err == nil {
				log.Debugf("pong/handleSubscriptionMessages: Open envelope: %v\n", msg)
				e.Messages <- msg
			}
		}
	}
}
