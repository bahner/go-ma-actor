package entity

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma-actor/p2p/pubsub"
	"github.com/bahner/go-ma/did"
	"github.com/bahner/go-ma/did/doc"
	"github.com/bahner/go-ma/msg"
	p2ppubsub "github.com/libp2p/go-libp2p-pubsub"
	log "github.com/sirupsen/logrus"
)

const (
	MESSAGES_BUFFERSIZE = 100
)

type Entity struct {

	// Live non-stored data
	// Context to be able to clean up entity.
	Ctx    context.Context
	Cancel context.CancelFunc
	Topic  *p2ppubsub.Topic

	//Stored data
	DID  did.DID
	Doc  *doc.Document
	Nick string // NOT the fragment. Can be but, nut needn't be.

	// Channels
	Messages chan *msg.Message
}

// Create a new entity from a DID
// In this case the DID is the string, not the struct.
func New(id string) (*Entity, error) {

	d, err := did.New(id)
	if err != nil {
		return nil, fmt.Errorf("entity/getorcreate: %w", err)
	}

	return NewFromDID(d)
}

// Create a new entity from a DID
func NewFromDID(d did.DID) (*Entity, error) {

	// Only 1 topic, but this is where it's at! One topic per entity.
	_topic, err := pubsub.GetOrCreateTopic(d.Id)
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

	err = e.getOrCreateAndSetNick()
	if err != nil {
		log.Debugf("entity/new: failed to get and set nick: %s", err)
	}

	// Cache the entity
	entities.Store(d.Id, e)

	return e, nil
}

// Get an entity from the global map.
func GetOrCreate(id string) (*Entity, error) {

	// Check if we have one cached
	e := load(GetDID(id)) // Be nice and see if this is a nick.
	if e != nil {
		err := e.getOrCreateAndSetNick() // If this fails nothing should happen.
		if err != nil {
			log.Debugf("GetOrCreate: failed to find existing nick: %s", err)
		}
		return e, nil
	}

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

	e, err := NewFromDID(id)
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

func Lookup(id string) (*Entity, error) {
	return GetOrCreate(GetDID(id))
}
