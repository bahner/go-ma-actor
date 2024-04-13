package entity

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma/api"
	"github.com/bahner/go-ma/did/doc"
	"github.com/fxamacker/cbor/v2"
	"github.com/ipfs/go-cid"
	log "github.com/sirupsen/logrus"
)

func (e *Entity) FetchDocument() (*doc.Document, error) {

	var document = &doc.Document{}

	c, err := api.ResolveRootCID(ip.String(), cached)
	if err != nil {
		return nil, cid.Cid{}, fmt.Errorf("fetchFromDID: %w", err)
	}

	log.Debugf("Fetching CID: %s", c)

	node, err := ipfsAPI.Dag().Get(context.Background(), c)
	if err != nil {
		return nil, cid.Cid{}, fmt.Errorf("fetchFromDID: %w", err)
	}

	err = cbor.Unmarshal(node.RawData(), document)
	if err != nil {
		return nil, cid.Cid{}, fmt.Errorf("fetchFromDID: %w", err)
	}

	err = document.Verify()
	if err != nil {
		return nil, cid.Cid{}, fmt.Errorf("fetchFromDID: %w", err)
	}

	log.Debugf("Fetched and cached document for : %s", d.Id)
	return document, c, nil

}

// Fetch the document and set it in the entity.
// If cached is true, the document will be fetched from the IPFS cache,
// if available.
func (e *Entity) FetchAndSetDocument() error {

	var (
		err error
	)

	if e.Doc == nil {
		// Fetch the document
		e.Doc, _, err = doc.FetchFromDID(e.DID, false)
		if err != nil {
			return err
		}
	}

	return nil
}
