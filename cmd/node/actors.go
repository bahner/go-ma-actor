package main

import (
	"fmt"
	"sync"

	"github.com/bahner/go-ma-actor/entity/actor"
	"github.com/bahner/go-ma/did/doc"
	log "github.com/sirupsen/logrus"
)

var actors *actorCache

type actorCache struct {
	store sync.Map
}

func init() {
	actors = new(actorCache)
}

// GetOrCreateEntity returns an entity from the cache or creates a new one
// The id is just the uniqgue name of the calling entity, not the full DID
func getOrCreateActor(id string) (*actor.Actor, error) {

	// Attempt to retrieve the entity from cache.
	// This is runtime, so entities will be generated at least once.
	if cachedEntity, ok := actors.Get(id); ok {
		if entity, ok := cachedEntity.(*actor.Actor); ok {
			log.Debugf("found topic: %s in entities cache.", id)
			return entity, nil // Successfully type-asserted and returned
		}
	}

	// Entity not found in cache, proceed to creation
	log.Debugf("getOrCreateEntity: GetOrCreateKeyset from vault: %s", id)
	k, err := getOrCreateKeysetFromVault(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get or create keyset: %w", err)
	}

	// Assuming entity.NewFromKeyset returns *actor.Actor
	a, err := actor.New(k)
	if err != nil {
		return nil, fmt.Errorf("failed to create entity: %w", err)
	}

	a.Entity.Doc, err = doc.NewFromKeyset(a.Keyset)
	if err != nil {
		return nil, fmt.Errorf("failed to create DID Document: %w", err)
	}

	_, err = a.Entity.Doc.Publish()
	if err != nil {
		return nil, fmt.Errorf("failed to publish DID Document: %w", err)
	}

	// Cache the newly created entity for future retrievals
	actors.Set(id, a)

	return a, nil
}

func (e *actorCache) Set(key string, value interface{}) {
	e.store.Store(key, value)
}

func (e *actorCache) Get(key string) (interface{}, bool) {
	return e.store.Load(key)
}
