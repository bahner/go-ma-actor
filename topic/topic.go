package topic

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma/msg"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

var topics = map[string]*Topic{}

const MESSAGES_BUFFERSIZE = 100

type Topic struct {
	Topic        *pubsub.Topic
	Subscription *pubsub.Subscription
	ctx          context.Context
	chDone       chan struct{}
	Messages     chan *msg.Message
}

func GetOrCreate(id string) (*Topic, error) {

	t := topics[id]
	if topics[id] != nil {
		return t, nil
	}

	pubsubTopic, err := getOrCreatePubSub(id)
	if err != nil {
		return nil, fmt.Errorf("failed to create topic: %v", err)
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
		return fmt.Errorf("topic/Close: failed to close topic: %v", err)
	}

	return nil
}

// Unsubscribe is used to stop the goroutine that is listening for messages.
func (t *Topic) Unsubscribe() error {

	close(t.chDone)

	return nil
}
