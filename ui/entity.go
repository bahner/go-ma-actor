package ui

import (
	"fmt"

	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma/did/doc"
)

func getOrCreateEntity(id string) (*entity.Entity, error) {

	var err error

	e := entity.GetOrCreate(id)

	// There should be a document there, but ...
	if e.Doc == nil {
		e.Doc, err = doc.GetOrFetch(id)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch DIDDOcument. %w", err)
		}
	}

	return e, nil
}
