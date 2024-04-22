package entity

import (
	"github.com/bahner/go-ma/did/doc"
)

func (e *Entity) FetchDocument() (*doc.Document, error) {

	d, _, err := doc.FetchFromDID(e.DID)
	if err != nil {
		return nil, err
	}

	return d, nil

}

// // Fetch the document and set it in the entity.
func (e *Entity) FetchAndSetDocument() error {

	var err error

	e.Doc, _, err = doc.FetchFromDID(e.DID)

	return err
}
