package space

import (
	"fmt"

	"github.com/bahner/go-ma-actor/actor"
	"github.com/bahner/go-ma-actor/topics"
	"github.com/bahner/go-ma/did/doc"
	"github.com/bahner/go-ma/msg"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

type Space struct {
	DID      string
	Document *doc.Document
	Private  *pubsub.Topic
	Public   *pubsub.Topic
	Messages chan *msg.Message
}

func Enter(id string, a *actor.Actor) (*Space, error) {

	if spaces[id] != nil {
		return spaces[id], nil
	}

	d, err := doc.Fetch(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch DID Document: %v", err)
	}

	keyAgreement, err := topics.GetOrCreate(id)
	if err != nil {
		return nil, fmt.Errorf("failed to join keyAgreement topic: %v", err)
	}

	assertionMethod, err := topics.GetOrCreate(d.AssertionMethod)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to assertionMethod topic: %v", err)
	}

	return &Space{
		DID:      id,
		Document: d,
		Private:  keyAgreement,
		Public:   assertionMethod,
	}, nil
}
