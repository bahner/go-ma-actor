package topic

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma/msg"
	"github.com/fxamacker/cbor/v2"
	log "github.com/sirupsen/logrus"
)

// Subscribe to a topic and start the subscription loop.
// NB! Not all entities can decrypt the envelopes, so this is not a method on the entity.
// Envelopes can be sent to the entity, but the entity must handle the decryption.
func (t *Topic) Subscribe(ctx context.Context, messages chan *msg.Message, envelopes chan *msg.Envelope) error {

	if messages == nil {
		return fmt.Errorf("messages channel must not be nil")
	}

	if envelopes == nil {
		log.Warnf("envelopes channel for topic subscription %s is nil. Ignoring envelopes.", t.Topic.String())
	}

	t.Ctx = ctx
	var err error

	t.Subscription, err = t.Topic.Subscribe()
	if err != nil {
		return fmt.Errorf("failed to subscribe to topic: %w", err)
	}

	go t.subscriptionLoop(messages, envelopes)

	return nil

}

func (t *Topic) subscriptionLoop(messages chan *msg.Message, envelopes chan *msg.Envelope) {
	for {
		select {
		case <-t.Ctx.Done():
			log.Debugf("Context cancelled, stopping envelope subscription loop for topic %s.", t.Topic.String())
			return
		case <-t.Done:
			log.Debugf("Channel done, stopping envelope subscription loop for topic %s.", t.Topic.String())
			return
		default:
			input, err := t.Subscription.Next(t.Ctx)
			if err != nil {
				log.Errorf("Error in envelope subscription loop: %v", err)
				continue
			}

			// First try to see if this fits a message, as that takes less time.
			var m *msg.Message

			err = cbor.Unmarshal(input.Data, &m)
			if err == nil {
				messages <- m
				continue
			}

			// If the envelopes channel is nil, we don't need to try to unmarshal envelopes.
			if envelopes == nil {
				var e *msg.Envelope
				err = cbor.Unmarshal(input.Data, &e)
				if err == nil {
					envelopes <- e
					continue
				}
			}

			log.Errorf("Failed to unmarshal message or envelope: %v", err)
		}
	}
}
