package entity

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma/did"
	"github.com/bahner/go-ma/did/doc"
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

	//Stored data
	DID did.DID
	Doc *doc.Document

	// Channels
	Messages chan *Message
}

// Creates a ned Entity from a DID and fetched the live document.
// This is used mostly for foreign entities.
func Init(d did.DID) (*Entity, error) {
	e, err := New(d)
	if err != nil {
		return nil, fmt.Errorf("entity.Fetch: %w", err)
	}

	err = e.fetchAndSetDocument()
	if err != nil {
		return nil, fmt.Errorf("entity.Fetch: %w", err)
	}

	err = e.joinTopic()
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

	ctx, cancel := context.WithCancel(context.Background())

	return &Entity{

		Ctx:    ctx,
		Cancel: cancel,

		DID: d,

		Messages: make(chan *Message, MESSAGES_BUFFERSIZE),
	}, nil
}
