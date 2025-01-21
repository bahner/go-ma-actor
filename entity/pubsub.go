package entity

// This is the asynchronous version of sending a message and as such
// does not guarantee delivery.

import (
	"fmt"

	"github.com/bahner/go-ma-actor/p2p/pubsub"
	"github.com/bahner/go-ma/did/doc"
	"github.com/bahner/go-ma/msg"
	log "github.com/sirupsen/logrus"
)

// Publishes a message or an envelope to the topic of the entity.
// As of now only p2p gossipsub is supported,
// but when other types of topics are supported, this function will be updated.
// An Entity can only have 1 topic defined at a time.
func (e *Entity) Publish(m msg.Msg) error {

	if e.Doc.Topic.ID == "" {
		return fmt.Errorf("entity.Publish: no topic set")
	}

	if e.Ctx == nil {
		return fmt.Errorf("entity.Publish: no context set")
	}

	topic, err := pubsub.GetOrCreateTopic(e.Doc.Topic.ID)
	if err != nil {
		return fmt.Errorf("entity.Publish: %w", err)
	}

	topic.Publish(e.Ctx, m.Bytes())

	return nil
}

func (e *Entity) joinTopic() error {

	var err error

	if e.Doc.Topic.Type != doc.DEFAULT_TOPIC_TYPE {
		log.Errorf("entity.Fetch: Topic of type %s not supported.", e.Doc.Topic.Type)
		return doc.ErrInvalidTopicType
	}

	// Join the typic, which caches the topic in the pubsub.
	_, err = pubsub.GetOrCreateTopic(e.Doc.Topic.ID)
	if err != nil {
		return fmt.Errorf("entity.Join: %w", err)
	}

	return nil
}
