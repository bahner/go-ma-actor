package actor

import (
	"fmt"

	"github.com/bahner/go-ma-actor/entity"
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

// // Get an entity from the global map.
// // The input is a full did string. If one is created it will have no Nick.
// // The function should do the required lookups to get the nick.
// // And verify the entity.
func GetOrCreate(id string) (*Actor, error) {

	// Creating a DID here implies validation before we try to load the actor.
	d, err := did.New(id)
	if err != nil {
		return nil, fmt.Errorf("actor.GetOrCreate: %w", err)
	}

	e := load(d.Id)
	if e != nil {
		return e, nil
	}

	k, err := set.GetOrCreate(d.Fragment)
	if err != nil {
		return nil, fmt.Errorf("entity/newfromdid: failed to get or create keyset: %w", err)
	}

	return New(k)
}
