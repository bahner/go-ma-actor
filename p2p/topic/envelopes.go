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
	e, err := envelope.UnmarshalFromCBOR(message.Data)
	if err != nil {
		log.Debugf("Failed to unmarshal envelope: %v", err)
	}

	return e, err

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
			log.Debugf("Context cancelled, stopping envelope subscription loop for topic %s.", t.Topic.String())
			return
		case <-t.chDone:
			log.Debugf("Channel done, stopping envelope subscription loop for topic %s.", t.Topic.String())
			return
		default:
			letter, err := t.NextEnvelope()
			if err != nil {
				log.Errorf("Error in envelope subscription loop: %v", err)
				continue
			}

			// If additional filtering or processing is required, do it here
			// For example, uncomment the following if you want to filter out messages sent by self
			// if letter.Sender == t.self {
			//     continue
			// }

			// Send valid messages onto the Envelopes channel
			t.Envelopes <- letter
		}
	}
}
