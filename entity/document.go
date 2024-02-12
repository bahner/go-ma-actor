package entity

import "github.com/bahner/go-ma/did/doc"

func (e *Entity) FetchDocument(cached bool) (*doc.Document, error) {

	var (
		err error
	)

	d := new(doc.Document)

	if e.Doc == nil {
		// Fetch the document
		d, err = doc.FetchFromDID(e.DID, cached)
		if err != nil {
			return d, err
		}
	}

	return e.Doc, nil
}

// Fetch the document and set it in the entity.
// If cached is true, the document will be fetched from the IPFS cache,
// if available.
func (e *Entity) FetchAndSetDocument(cached bool) error {

	var (
		err error
	)

	if e.Doc == nil {
		// Fetch the document
		e.Doc, err = doc.FetchFromDID(e.DID, cached)
		if err != nil {
			return err
		}
	}

	return nil
}
