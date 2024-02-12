package actor

import (
	"fmt"

	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma/did"
	"github.com/bahner/go-ma/key/set"
	"github.com/bahner/go-ma/msg"
)

const ENVELOPES_BUFFERSIZE = 100

// An actor is just an entity with a keyset.
type Actor struct {
	Entity *entity.Entity

	Keyset *set.Keyset

	Envelopes chan *msg.Envelope
}

// Create a new entity from a DID and a Keyset. We need both.
// The DID is to verify the entity, and the keyset is to create the
// DID Document.
func New(d *did.DID, k *set.Keyset) (*Actor, error) {

	if k == nil {
		return nil, fmt.Errorf("entity/new: no keyset")
	}

	// Hrm. Use the entity context or create own ...
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	e, err := entity.New(d)
	if err != nil {
		return nil, err
	}

	a := &Actor{
		Entity: e,
		Keyset: k,
	}

	a.CreateDocument(d.String())

	// Cache the entity
	store(a)

	return a, nil
}

// Create a new entity from a DID and use fragment as nick.
func NewFromDID(id string, nick string) (*Actor, error) {

	d, err := did.New(id)
	if err != nil {
		return nil, fmt.Errorf("entity/newfromdid: failed to create did from ipnsKey: %w", err)
	}

	k, err := set.GetOrCreate(d.Fragment)
	if err != nil {
		return nil, fmt.Errorf("entity/newfromdid: failed to get or create keyset: %w", err)
	}

	return New(d, k)
}

// // Get an entity from the global map.
// // The input is a full did string. If one is created it will have no Nick.
// // The function should do the required lookups to get the nick.
// // And verify the entity.
func GetOrCreate(id string) (*Actor, error) {

	if id == "" {
		return nil, fmt.Errorf("entity/getorcreate: empty id")
	}

	if !did.IsValidDID(id) {
		return nil, fmt.Errorf("entity/getorcreate: invalid id")
	}

	var err error

	e := load(id)
	if e != nil {
		return e, nil
	}

	e, err = NewFromDID(id, "")
	if err != nil {
		return nil, fmt.Errorf("entity/getorcreate: failed to create entity: %w", err)
	}

	err = e.Verify()
	if err != nil {
		return nil, fmt.Errorf("entity/getorcreate: failed to verify created entity: %w", err)
	}

	return e, nil
}
