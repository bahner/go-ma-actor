package entity

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma-actor/alias"
	"github.com/bahner/go-ma-actor/p2p/topic"
	"github.com/bahner/go-ma/did"
	"github.com/bahner/go-ma/did/doc"
	"github.com/bahner/go-ma/key/set"
	"github.com/bahner/go-ma/msg"
)

const (
	MESSAGES_BUFFERSIZE  = 100
	ENVELOPES_BUFFERSIZE = 100
)

type Entity struct {
	Ctx        context.Context
	CancelFunc context.CancelFunc

	// External structs
	DID   *did.DID
	Doc   *doc.Document
	Topic *topic.Topic

	// Only keyset maybe nil
	Keyset *set.Keyset

	// Channel for incoming messages
	Messages  chan *msg.Message
	Envelopes chan *msg.Envelope

	// Nick is pretty much the same as the fragment of the DID
	// But You set this, so you can trust it.
	Nick string
}

// Create a new entity from a DID and give it a nick.
func New(d *did.DID, k *set.Keyset, nick string) (*Entity, error) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_topic, err := topic.GetOrCreate(d.String())
	if err != nil {
		return nil, fmt.Errorf("entity/new: failed to create new topic: %w", err)
	}

	_doc, err := doc.Fetch(d.String(), true) // Accept cached version
	if err != nil {
		return nil, fmt.Errorf("entity/new: failed to create new document: %w", err)
	}

	// Look up nick if not set else set it.
	if nick == "" {
		nick = alias.LookupEntityDID(d.String())
	}

	e := &Entity{
		Ctx:        ctx,
		CancelFunc: cancel,

		Nick: nick,

		DID:   d,
		Doc:   _doc,
		Topic: _topic,

		Keyset: k,

		Messages: make(chan *msg.Message, MESSAGES_BUFFERSIZE),
	}

	add(e)

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
func GetOrCreate(id string) (*Entity, error) {

	var err error

	e := get(id)
	if e != nil {
		e.Nick = alias.LookupEntityDID(id)
		return e, nil
	}

	e, err = NewFromDID(id, "")
	if err != nil {
		return nil, fmt.Errorf("entity/getorcreate: failed to create entity: %w", err)
	}

	return e, nil
}

func (e *Entity) Leave() {
	e.CancelFunc()
}

// Takes a message channel and and actor entity and recieves messages
// The actor is required to decrypt the envelopes.
func (e *Entity) Enter(actor *Entity) error {

	err := actor.Verify()
	if err != nil {
		return fmt.Errorf("entity/start: failed to verify actor: %w", err)
	}

	if actor.Keyset == nil {
		return fmt.Errorf("entity/start: actor has no keyset")
	}

	go e.receiveMessages(actor)

	return nil
}
