package entity

import (
	"fmt"

	"github.com/bahner/go-ma-actor/alias"
	"github.com/bahner/go-ma-actor/p2p/pubsub"
	"github.com/bahner/go-ma/did"
	"github.com/bahner/go-ma/did/doc"
	"github.com/bahner/go-ma/key/set"
	"github.com/bahner/go-ma/msg"
	p2ppubsub "github.com/libp2p/go-libp2p-pubsub"
)

const (
	MESSAGES_BUFFERSIZE  = 100
	ENVELOPES_BUFFERSIZE = 100
)

type Entity struct {

	// External structs
	DID *did.DID
	Doc *doc.Document

	Topic        *p2ppubsub.Topic
	Subscription *Subscription

	// Only keyset maybe nil
	Keyset *set.Keyset

	// Channels
	Messages  chan *msg.Message
	Envelopes chan *msg.Envelope

	// Nick is pretty much the same as the fragment of the DID
	// But You set this, so you can trust it.
	Nick string
}

// Create a new entity from a DID and give it a nick.
func New(d *did.DID, k *set.Keyset, nick string) (*Entity, error) {

	ps := pubsub.Get()

	// Only 1 topic, but this is where it's at! One topic er entity.
	_topic, err := ps.Join(d.String())
	if err != nil {
		return nil, fmt.Errorf("entity/new: failed to join topic: %w", err)
	}

	_doc, err := doc.Fetch(d.String(), true) // Accept cached version
	if err != nil {
		return nil, fmt.Errorf("entity/new: failed to create new document: %w", err)
	}

	// Look up nick if not set else set it.
	if nick == "" {
		nick = alias.GetOrCreateEntityAlias(d.String())
	}

	e := &Entity{

		Nick: nick,

		DID:   d,
		Doc:   _doc,
		Topic: _topic,

		Messages:  make(chan *msg.Message, MESSAGES_BUFFERSIZE),
		Envelopes: make(chan *msg.Envelope, ENVELOPES_BUFFERSIZE),

		Keyset: k,
	}

	// Cache the entity
	cache(e)
	e.Subscription, err = e.Subscribe()
	if err != nil {
		return nil, fmt.Errorf("entity/new: failed to subscribe to topic: %w", err)
	}

	return e, nil
}

// Create a new entity from a DID and use fragment as nick.
func NewFromDID(id string, nick string) (*Entity, error) {

	d, err := did.New(id)
	if err != nil {
		return nil, fmt.Errorf("entity/newfromdid: failed to create did from ipnsKey: %w", err)
	}

	return New(d, nil, nick)
}

// Get an entity from the global map.
// The input is a full did string. If one is created it will have no Nick.
// The function should do the required lookups to get the nick.
// And verify the entity.
func GetOrCreate(id string) (*Entity, error) {

	if id == "" {
		return nil, fmt.Errorf("entity/getorcreate: empty id")
	}

	if !did.IsValidDID(id) {
		return nil, fmt.Errorf("entity/getorcreate: invalid id")
	}

	var err error

	e := get(id)
	if e != nil {
		e.Nick = alias.GetOrCreateEntityAlias(id)
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
