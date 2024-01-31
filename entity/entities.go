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

// List returns a list of entities
func List() []*Entity {
	var result []*Entity
	for _, p := range entities {
		result = append(result, p)
	}
	return result
}

// Remove removes a entity from the map
func Delete(id string) {
	delete(entities, id)
}

func ListAliases() []string {
	var result []string
	for _, e := range entities {
		result = append(result, e.Nick)
	}
	return result
}

func ListDIDs() []string {
	var result []string
	for _, e := range entities {
		result = append(result, e.DID.String())
	}
	return result
}
