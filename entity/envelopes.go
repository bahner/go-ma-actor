package entity

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma/msg"
	"github.com/fxamacker/cbor/v2"
	log "github.com/sirupsen/logrus"
)

func (e *Entity) Subscribe(ctx context.Context) error {

	e.Topic.Ctx = ctx
	var err error

	e.Topic.Subscription, err = e.Topic.Topic.Subscribe()
	if err != nil {
		return fmt.Errorf("failed to subscribe to topic: %w", err)
	}

	go e.subscriptionLoop()

	return nil

}

// The subscription lop must run on the entity,
// as the private key is needed to open the envelope.
func (e *Entity) subscriptionLoop() {

	go e.envelopeHandler()

	for {
		select {
		case <-e.Topic.Ctx.Done():
			log.Debugf("Context cancelled, stopping envelope subscription loop for entity topic %s.", e.Topic.Topic.String())
			return
		case <-e.Topic.Done:
			log.Debugf("Channel done, stopping envelope subscription loop for entity topic %s.", e.Topic.Topic.String())
			return
		default:
			input, err := e.Topic.Subscription.Next(e.Topic.Ctx)
			if err != nil {
				log.Errorf("Error in envelope subscription loop: %v", err)
				continue
			}

			// First try to see if this fits a message, as that takes less time.
			var m *msg.Message

			err = cbor.Unmarshal(input.Data, &m)
			if err == nil {
				err := m.Verify()
				if err == nil {
					e.Messages <- m
				} else {
					log.Errorf("Failed to verify message: %v", err)
				}
				continue
			} // Nothing to log here. Maybe debug, but nothing to worry about.

			// Then go for the jugular
			var me *msg.Envelope

			err = cbor.Unmarshal(input.Data, &me)
			if err == nil {
				e.Envelopes <- me
				continue
			}

			log.Errorf("Failed to unmarshal message or envelope: %v", err)
		}
	}
}

func (e *Entity) envelopeHandler() {

	envelope := <-e.Envelopes

	for {
		select {
		case <-e.Topic.Done:
			return
		default:
			if msg, err := envelope.Open(e.Keyset.EncryptionKey.PrivKey[:]); err == nil {
				e.Messages <- msg
			}
		}
	}
}
