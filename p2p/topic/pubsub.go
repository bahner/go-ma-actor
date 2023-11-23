package topic

import (
	"github.com/bahner/go-ma-actor/p2p/pubsub"
	p2ppubsub "github.com/libp2p/go-libp2p-pubsub"
)

var pubsubTopics = map[string]*p2ppubsub.Topic{}

func getOrCreatePubSub(id string) (*p2ppubsub.Topic, error) {

	t := getPubSub(id)
	if t != nil {
		return t, nil
	}

	return createPubSub(id)
}

func createPubSub(id string) (*p2ppubsub.Topic, error) {

	ps := pubsub.Get()

	pst, err := ps.Join(id)
	if err != nil {
		return nil, err
	}

	addPubSub(id, pst)

	return pst, nil
}

func addPubSub(id string, topic *p2ppubsub.Topic) {
	pubsubTopics[id] = topic
}

func getPubSub(id string) *p2ppubsub.Topic {
	return pubsubTopics[id]
}

func deletePubSub(id string) {
	delete(pubsubTopics, id)
}
