package topic

import (
	"github.com/bahner/go-ma-actor/p2p/pubsub"
	p2ppubsub "github.com/libp2p/go-libp2p-pubsub"
)

var pubsubTopics = map[string]*p2ppubsub.Topic{}

func getOrCreatePubSub(id string) (*p2ppubsub.Topic, error) {

	t := pubsubTopics[id]
	if t != nil {
		return t, nil
	}

	return createPubSub(id)
}

func createPubSub(id string) (*p2ppubsub.Topic, error) {

	ps := pubsub.Get()

	t, err := ps.Join(id)
	if err != nil {
		return nil, err
	}

	return t, nil
}
