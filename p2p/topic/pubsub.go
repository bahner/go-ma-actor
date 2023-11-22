package topic

import (
	"github.com/bahner/go-ma-actor/config"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

var pubsubTopics = map[string]*pubsub.Topic{}

func getOrCreatePubSub(id string) (*pubsub.Topic, error) {

	t := pubsubTopics[id]
	if t != nil {
		return t, nil
	}

	return createPubSub(id)
}

func createPubSub(id string) (*pubsub.Topic, error) {

	ps := config.GetPubSub()

	t, err := ps.Join(id)
	if err != nil {
		return nil, err
	}

	return t, nil
}
