package topics

import (
	"github.com/bahner/go-ma-actor/config"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

var topics = map[string]*pubsub.Topic{}

func GetOrCreate(id string) (*pubsub.Topic, error) {

	t := topics[id]
	if t != nil {
		return t, nil
	}

	return Create(id)
}

func Create(id string) (*pubsub.Topic, error) {

	ps := config.GetPubSub()

	t, err := ps.Join(id)
	if err != nil {
		return nil, err
	}

	topics[id] = t

	return t, nil
}

// Fetch a topic if it is known.
func Get(id string) *pubsub.Topic {
	return topics[id]
}

// Delete a topic if it is known.
func Delete(id string) error {

	t := topics[id]
	if t == nil {
		return nil
	}

	t.Close()

	delete(topics, id)

	return nil
}

func CloseAll() {
	for _, t := range topics {
		t.Close()
	}
}

func List() []string {
	var list []string
	for id := range topics {
		list = append(list, id)
	}
	return list
}

func Count() int {
	return len(topics)
}

func Clear() {
	topics = map[string]*pubsub.Topic{}
}

func Exists(id string) bool {
	return topics[id] != nil
}

func IsEmpty() bool {
	return len(topics) == 0
}

func IsNotEmpty() bool {
	return len(topics) > 0
}
