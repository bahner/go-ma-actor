package topic

import (
	"context"
	"encoding/hex"

	"github.com/bahner/go-ma/msg"
	log "github.com/sirupsen/logrus"
	"lukechampine.com/blake3"
)

func (t *Topic) SubscribeEnvelopes(ctx context.Context) (envelopes <-chan *msg.Envelope) {

	t.ctx = ctx
	var err error

	t.Subscription, err = t.Topic.Subscribe()
	if err != nil {
		log.Errorf("Failed to subscribe to topic: %v", err)
		return
	}

	t.Envelopes = make(chan *msg.Envelope, ENVELOPES_BUFFERSIZE)

	go t.envelopeSubscriptionLoop()

	return t.Envelopes

}

func (t *Topic) NextEnvelope() (*msg.Envelope, error) {

	message, err := t.Subscription.Next(t.ctx)
	if err != nil {
		return nil, err
	}

	if log.GetLevel() >= log.DebugLevel {
		// blake3 is very fast, so this is not a problem in debugging mode
		// This is just so we can see that messages are actually received
		bs := blake3.Sum256(message.Data)
		checksum := hex.EncodeToString(bs[:])
		log.Debugf("Received message with checksum: %s", checksum)
	}

	// Here we should distinguish between packed and unpacked envelopes
	e, err := msg.UnmarshalEnvelopeFromCBOR(message.Data)
	if err != nil {
		log.Errorf("Failed to unmarshal envelope: %v", err)
	}

	return e, err

}

// Publish a message to the topic.
// NB! Check that it's the correct topic!
func (t *Topic) SendMessage(m *msg.Message) error {

	// Should this just be a goroutine. Message delivery is not guaranteed anyway.
	return m.Send(t.ctx, t.Topic)

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
