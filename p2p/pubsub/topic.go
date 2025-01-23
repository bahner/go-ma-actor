package pubsub

import (
	"fmt"
	"sync"

	p2ppubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
)

var topic sync.Map

func GetOrCreateTopic(topicName string) (*p2ppubsub.Topic, error) {
	t, ok := topic.Load(topicName)
	if !ok {
		ps := Get()
		var err error
		t, err = ps.Join(topicName)
		if err != nil {
			return nil, fmt.Errorf("entity/getOrCreateTopic: failed to join topic: %w", err)
		}
		topic.Store(topicName, t)
	}

	return t.(*p2ppubsub.Topic), nil
}

// List the peers for a gossipsub topic.
// Returns an empty slice if the topic does not exist.
func ListPeers(topicName string) []peer.ID {
	t, err := GetOrCreateTopic(topicName)
	if err != nil {
		return []peer.ID{}
	}

	return t.ListPeers()
}
