package topic

import (
	"context"
	"fmt"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

var (
	err    error
	topics = map[string]*Topic{}
)

type Topic struct {
	Ctx context.Context

	Done chan struct{}

	Topic        *pubsub.Topic
	Subscription *pubsub.Subscription
}

func GetOrCreate(id string) (*Topic, error) {

	t, exists := Get(id)
	if exists {
		return t, nil
	}

	t = &Topic{
		Done: make(chan struct{}),
	}

	// Topic
	t.Topic, err = getOrCreatePubSub(id)
	if err != nil {
		return nil, fmt.Errorf("failed to create topic: %w", err)
	}
	addPubSub(id, t.Topic)

	// Subscription
	t.Subscription, err = t.Topic.Subscribe()
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to topic: %w", err)
	}

	// Add it to the list of topics
	add(t)
	return t, nil
}

func Get(id string) (*Topic, bool) {

	t, exists := topics[id]

	return t, exists
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

	close(t.Done)

	return nil
}

func (t *Topic) Delete() {

	deletePubSub(t.Topic.String())
	delete(topics, t.Topic.String())
}

func add(t *Topic) {
	topics[t.Topic.String()] = t
}
