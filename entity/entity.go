package entity

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma/did"
	"github.com/bahner/go-ma/did/doc"
	"github.com/bahner/go-ma/msg"
	p2ppubsub "github.com/libp2p/go-libp2p-pubsub"
)

const (
	MESSAGES_BUFFERSIZE = 100
)

type Entity struct {
	// Context to be able to clean up entity.
	Ctx    context.Context
	Cancel context.CancelFunc

	// External structs
	DID did.DID
	Doc *doc.Document

	Topic *p2ppubsub.Topic

	// Channels
	Messages chan *msg.Message
}

// Create a new entity from a DID and give it a nick.
func New(d did.DID) (*Entity, error) {

	// Only 1 topic, but this is where it's at! One topic per entity.
	_topic, err := getOrCreateTopic(d.Id)
	if err != nil {
		return nil, fmt.Errorf("entity/new: failed to join topic: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	e := &Entity{

		Ctx:    ctx,
		Cancel: cancel,

		DID:   d,
		Topic: _topic,

		Messages: make(chan *msg.Message, MESSAGES_BUFFERSIZE),
	}

	// Cache the entity
	store(e)

	return e, nil
}

// Create a new entity from a DID.
// In this case the DID is the strinf, not the struct.
func NewFromDID(id string) (*Entity, error) {

	d, err := did.New(id)
	if err != nil {
		return nil, fmt.Errorf("entity/getorcreate: %w", err)
	}

	return New(d)
}

// Get an entity from the global map.
// The input is a full did string. If one is created it will have no Nick.
// The function should do the required lookups to get the nick.
// And verify the entity.
func GetOrCreate(id string) (*Entity, error) {

	// Check if we have one cahced
	e := load(id)
	if e != nil {
		return e, nil
	}

	// Create a DID from the string
	d, err := did.New(id)
	if err != nil {
		return nil, fmt.Errorf("GetOrCreate: %w", err)
	}

	return GetOrCreateFromDID(d)

}

// Get an entity from the global map.
// The input is a full did string. If one is created it will have no Nick.
// The function should do the required lookups to get the nick.
// And verify the entity.
func GetOrCreateFromDID(id did.DID) (*Entity, error) {

	e, err := New(id)
	if err != nil {
		return nil, fmt.Errorf("GetOrCreateFromDID: %w", err)
	}

	// Fetch the document
	err = e.FetchAndSetDocument(false)
	if err != nil {
		return nil, fmt.Errorf("GetOrCreateFromDID: %w", err)
	}

	err = e.Verify()
	if err != nil {
		return nil, fmt.Errorf("GetOrCreateFromDID: %w", err)
	}

	return e, nil
}
