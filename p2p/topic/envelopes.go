package topic

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma/msg/envelope"
	log "github.com/sirupsen/logrus"
)

func (t *Topic) SubscribeEnvelopes(ctx context.Context) (envelopes <-chan *envelope.Envelope) {

	t.ctx = ctx
	var err error

	t.Subscription, err = t.Topic.Subscribe()
	if err != nil {
		log.Errorf("Failed to subscribe to topic: %v", err)
		return
	}

	t.Envelopes = make(chan *envelope.Envelope, ENVELOPES_BUFFERSIZE)

	go t.envelopeSubscriptionLoop()

	return t.Envelopes

}

func (t *Topic) NextEnvelope() (*envelope.Envelope, error) {

	message, err := t.Subscription.Next(t.ctx)
	if err != nil {
		return nil, err
	}

	// Here we should distinguish between packed and unpacked envelopes
	return envelope.UnmarshalFromCBOR(message.Data)

}

// Publish a message to the topic.
// NB! Check that it's the correct topic!
func (t *Topic) SendEnvelope(e *envelope.Envelope) error {

	data, err := e.MarshalToCBOR()
	if err != nil {
		return err
	}

	err = t.Topic.Publish(t.ctx, data)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

func (t *Topic) envelopeSubscriptionLoop() {
	for {
		select {
		case <-t.ctx.Done():
			return
		case <-t.chDone:
			return
		default:
			letter, err := t.NextEnvelope()
			if err != nil {
				close(t.Envelopes)
				return
			}

			// See everything for now.
			// // only forward messages delivered by others
			// if msg.ReceivedFrom == cr.self {
			// 	continue
			// }

			// send valid messages onto the Messages channel
			t.Envelopes <- letter

		}
	}
}
