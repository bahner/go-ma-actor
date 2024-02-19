package actor

import (
	"sync"
)

var actors sync.Map

// store adds an entity to the map
func store(a *Actor) {
	actors.Store(a.Entity.DID.Id, a)
}

// load returns an entity from the map
func load(id string) *Actor {
	if entity, ok := actors.Load(id); ok {
		return entity.(*Actor)
	}
	return nil
}
