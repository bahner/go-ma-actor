package topic

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma/msg"
	"github.com/bahner/go-ma/msg/envelope"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

var topics = map[string]*Topic{}

const (
	MESSAGES_BUFFERSIZE  = 100
	ENVELOPES_BUFFERSIZE = 100
)

type Topic struct {
	ctx context.Context

	chDone chan struct{}

	Messages  chan *msg.Message
	Envelopes chan *envelope.Envelope

	Topic        *pubsub.Topic
	Subscription *pubsub.Subscription
}

func GetOrCreate(id string) (*Topic, error) {

	t := topics[id]
	if topics[id] != nil {
		return t, nil
	}

	pubsubTopic, err := getOrCreatePubSub(id)
	if err != nil {
		return nil, fmt.Errorf("failed to create topic: %w", err)
	}

	return &Topic{
		Topic:  pubsubTopic,
		chDone: make(chan struct{}),
	}, nil
}

// Close a topic if it is known.
func (t *Topic) Close() error {

	t.Unsubscribe()

	if t.Topic == nil {
		return nil
	}

	err := t.Topic.Close()
	if err != nil {
		return fmt.Errorf("topic/Close: failed to close topic: %w", err)
	}

	return nil
}

// Unsubscribe is used to stop the goroutine that is listening for messages.
func (t *Topic) Unsubscribe() error {

	close(t.chDone)

	return nil
}
