package entity

import (
	"github.com/bahner/go-ma/did/doc"
)

// Fetch the document and set it in the entity.
func (e *Entity) fetchAndSetDocument() error {

	var err error

	e.Doc, _, err = doc.FetchFromDID(e.DID)

	return err
}
