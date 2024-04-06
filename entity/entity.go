package entity

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma-actor/p2p/pubsub"
	"github.com/bahner/go-ma/did"
	"github.com/bahner/go-ma/did/doc"
	p2ppubsub "github.com/libp2p/go-libp2p-pubsub"
)

const (
	MESSAGES_BUFFERSIZE = 100
	entityPrefix        = "entity:"
)

type Entity struct {

	// Live non-stored data
	// Context to be able to clean up entity.
	Ctx    context.Context
	Cancel context.CancelFunc
	Topic  *p2ppubsub.Topic

	//Stored data
	DID did.DID
	Doc *doc.Document

	// Channels
	Messages chan *Message
}

// Creates a ned Entity from a DID and fetched the live document.
// This is used mostly for foreign entities.
func Fetch(d did.DID) (*Entity, error) {
	e, err := New(d)
	if err != nil {
		return nil, fmt.Errorf("entity.Fetch: %w", err)
	}

	err = e.FetchAndSetDocument()
	if err != nil {
		return nil, fmt.Errorf("entity.Fetch: %w", err)
	}

	err = e.Verify()
	if err != nil {
		return nil, fmt.Errorf("entity.Fetch: %w", err)
	}

	return e, nil

}

// Creates a new Entity from a DID, but does not fetch the document.
// Use this when creating a new actor for instance
func New(d did.DID) (*Entity, error) {

	// Only 1 topic, but this is where it's at! One topic per entity.
	_topic, err := pubsub.GetOrCreateTopic(d.Id)
	if err != nil {
		return nil, fmt.Errorf("entity/new: failed to join topic: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return &Entity{

		Ctx:    ctx,
		Cancel: cancel,

		DID:   d,
		Topic: _topic,

		Messages: make(chan *Message, MESSAGES_BUFFERSIZE),
	}, nil
}
