package actor

import (
	"fmt"

	"github.com/bahner/go-ma-actor/db"
	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma-actor/p2p"
	"github.com/bahner/go-ma/did"
	"github.com/bahner/go-ma/did/doc"
	"github.com/bahner/go-ma/key/set"
	"github.com/bahner/go-ma/msg"
)

const ENVELOPES_BUFFERSIZE = 100

// An actor is just an entity with a keyset.
type Actor struct {
	// The entity is the core of the actor.
	Entity *entity.Entity

	// The keyset is used to sign messages.
	Keyset set.Keyset

	// P2P The actors libp2p host.
	P2P *p2p.P2P

	// Envelopes are here. Messages comes to the entity.
	// So Entity.Messages
	Envelopes chan *msg.Envelope

	// Location of the actor.
	Location *entity.Entity

	// A handler for the dot messages.
	MessageHandler func(*msg.Message) error
}

// Create a new entity from a DID and a Keyset. We need both.
// The DID is to verify the entity, and the keyset is to create the
// DID Document.
func New(k set.Keyset) (*Actor, error) {

	err := k.Verify()
	if err != nil {
		return nil, fmt.Errorf("entity/new: failed to verify keyset: %w", err)
	}

	e, err := entity.New(k.DID)
	if err != nil {
		return nil, err
	}

	a := &Actor{
		Entity:    e,
		Keyset:    k,
		Envelopes: make(chan *msg.Envelope, ENVELOPES_BUFFERSIZE),
	}

	// Set a default message handler
	a.MessageHandler = a.defaultMessageHandler

	// Set and publish the actor DID Document
	a.Entity.Doc, err = doc.NewFromKeyset(a.Keyset)
	if err != nil {
		panic(err)
	}

	store(a)

	return a, nil
}

// Takes a DID String as input and returns the actor.
func GetOrCreate(id string) (*Actor, error) {

	a := load(id)
	if a != nil {
		return a, nil
	}

	d, err := did.NewFromString(id)
	if err != nil {
		return nil, fmt.Errorf("actor.GetOrCreate: failed to create DID: %w", err)
	}

	identity, err := db.GetOrCreateIdentity(d.Fragment)
	if err != nil {
		return nil, fmt.Errorf("actor.GetOrCreate: %w", err)
	}

	k, err := set.New(identity, d.Fragment)
	if err != nil {
		return nil, fmt.Errorf("entity/newfromdid: failed to get or create keyset: %w", err)
	}

	return New(k)
}
