package entity

import (
	"fmt"
	"sync"

	"github.com/bahner/go-ma-actor/p2p/pubsub"
	p2ppubsub "github.com/libp2p/go-libp2p-pubsub"
)

var topic sync.Map

func getOrCreateTopic(topicName string) (*p2ppubsub.Topic, error) {
	t, ok := topic.Load(topicName)
	if !ok {
		ps := pubsub.Get()
		var err error
		t, err = ps.Join(topicName)
		if err != nil {
			return nil, fmt.Errorf("entity/getOrCreateTopic: failed to join topic: %w", err)
		}
		topic.Store(topicName, t)
	}

	return t.(*p2ppubsub.Topic), nil
}
