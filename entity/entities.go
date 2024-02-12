package entity

import (
	"sync"
)

var entities sync.Map

// store adds an entity to the map
func store(e *Entity) {
	entities.Store(e.DID.String(), e)
}

// load returns an entity from the map
func load(id string) *Entity {
	if entity, ok := entities.Load(id); ok {
		return entity.(*Entity) // Type assert to *Entity, since Load returns an interface{}
	}
	return nil
}
