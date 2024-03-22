package entity

import (
	"context"
	"fmt"
	"strings"

	"github.com/bahner/go-ma-actor/db"
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

	return &Entity{

		Ctx:    ctx,
		Cancel: cancel,

		DID:   d,
		Topic: _topic,

		Messages: make(chan *Message, MESSAGES_BUFFERSIZE),
	}, nil
}

func GetOrCreate(id string, cached bool) (*Entity, error) {

	d, err := did.New(id)
	if err != nil {
		return nil, fmt.Errorf("entity/getorcreate: %w", err)
	}

	return GetOrCreateFromDID(d, cached)
}

// Sets a node in the database
// takes new did and nick. If an old  did  for the alias exists it is removed.
// This makes this the only alias for the DID and the only complex function in this file.
func (e Entity) SetNick(nick string) error {

	prefixBytes := []byte(entityNickPrefix)
	nickBytes := []byte(nick)
	idBytes := []byte(e.DID.Id)

	return db.Upsert(prefixBytes, nickBytes, idBytes)

}

// Returns the Entitty's nick. If it doesn't exist it returns the DID.
func (e Entity) Nick() string {

	idBytes := []byte(e.DID.Id)

	key, err := db.Lookup(idBytes)
	if err != nil {
		return e.DID.Id
	}

	return strings.TrimPrefix(string(key), entityNickPrefix)
}

// Get an entity from the global map.
// The input is a full did string. If one is created it will have no Nick.
// The function should do the required lookups to get the nick.
// And verify the entity.
// Cached is whether or not to use the cached copy of the entity document in IPFS
func GetOrCreateFromDID(id did.DID, cached bool) (*Entity, error) {

	e, err := NewFromDID(id)
	if err != nil {
		return nil, fmt.Errorf("GetOrCreateFromDID: %w", err)
	}

	// Fetch the document
	err = e.FetchAndSetDocument(cached)
	if err != nil {
		return nil, fmt.Errorf("GetOrCreateFromDID: %w", err)
	}

	err = e.Verify()
	if err != nil {
		return nil, fmt.Errorf("GetOrCreateFromDID: %w", err)
	}

	return e, nil
}
