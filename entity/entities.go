package entity

import (
	"errors"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var (
	ErrNotFound = errors.New("Entity not found")
	entities    sync.Map
)

// load returns an entity from the map
func load(id string) *Entity {
	if entity, ok := entities.Load(id); ok {
		return entity.(*Entity) // Type assert to *Entity, since Load returns an interface{}
	}
	return nil
}
