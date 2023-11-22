package topic

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma/msg"
	log "github.com/sirupsen/logrus"
)

func (t *Topic) SubscribeMessages(ctx context.Context) (messages <-chan *msg.Message) {

	t.ctx = ctx
	var err error

	t.Subscription, err = t.Topic.Subscribe()
	if err != nil {
		log.Errorf("Failed to subscribe to topic: %v", err)
		return
	}

	t.Messages = make(chan *msg.Message, MESSAGES_BUFFERSIZE)

	go t.messageSubscriptionLoop()

	return t.Messages

}

func (t *Topic) nextMessage() (*msg.Message, error) {

	message, err := t.Subscription.Next(t.ctx)
	if err != nil {
		return nil, err
	}

	// Here we should distinguish between different message types
	unpackedMessage, err := msg.Unpack(string(message.Data))
	if err != nil {
		return nil, fmt.Errorf("failed to unpack message: %v", err)
	}

	return unpackedMessage, nil
}

// Publish a message to the topic.
// NB! Check that it's the correct topic!
func (t *Topic) PublishMessage(m *msg.Message) error {

	data, err := m.Pack()
	if err != nil {
		return err
	}

	err = t.Topic.Publish(t.ctx, []byte(data))
	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}

	return nil
}

func (t *Topic) messageSubscriptionLoop() {
	for {
		select {
		case <-t.ctx.Done():
			return
		case <-t.chDone:
			return
		default:
			message, err := t.nextMessage()
			if err != nil {
				close(t.Messages)
				return
			}

			// See everything for now.
			// // only forward messages delivered by others
			// if msg.ReceivedFrom == cr.self {
			// 	continue
			// }

			// send valid messages onto the Messages channel
			t.Messages <- message

		}
	}
}
