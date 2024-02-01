package entity

var entities map[string]*Entity

func init() {
	entities = make(map[string]*Entity)
}

// add adds a entity to the map
func add(e *Entity) {
	entities[e.DID.String()] = e
}

// get returns a entity from the map
func get(id string) *Entity {
	return entities[id]
}
