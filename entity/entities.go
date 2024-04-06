package entity

import (
	"fmt"
	"sync"

	"github.com/bahner/go-ma/did"
)

var entities *sync.Map

func init() {
	entities = new(sync.Map)
}

// GetOrCreate checks if an Entity with the given DID string exists and returns it;
// otherwise, it creates a new Entity, stores it, and returns it.
// The boolean signifies whether to use a cached Entity document in IPFS, or
// to perform a network search. This might take time at the scale of minutes.
func GetOrCreate(didString string) (*Entity, error) {

	d, err := did.New(didString)
	if err != nil {
		return nil, fmt.Errorf("entity/getorcreate: %w", err)
	}

	if entity, ok := entities.Load(d.Id); ok {
		return entity.(*Entity), nil
	}

	entity, err := New(d)
	if err != nil {
		return nil, err
	}

	entities.Store(entity.DID.Id, entity)

	return entity, nil
}
