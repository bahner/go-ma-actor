package entity

import (
	"fmt"

	"github.com/bahner/go-ma/did"
	"github.com/bahner/go-ma/did/doc"
)

var err error

type Entity struct {
	// ID is the entity's ID
	DID string
	Doc *doc.Document
	// Name is the entity's name
	Alias string
}

func New(id string, alias string) (*Entity, error) {

	e := &Entity{
		DID:   id,
		Alias: alias,
	}

	e.Doc, err = doc.FetchFromDID(id)
	if err != nil {
		return nil, fmt.Errorf("entity/newfromdid: failed to fetch document: %w", err)
	}

	Add(e)

	return e, nil
}

func NewFromDID(id string) (*Entity, error) {

	alias := did.GetFragment(id)

	return New(id, alias)
}

func (e *Entity) IsValid() bool {
	return did.IsValidDID(e.DID)
}

func (e *Entity) GetDoc() *doc.Document {
	return e.Doc
}
