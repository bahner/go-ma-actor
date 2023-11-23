package entity

var entities map[string]*Entity

func init() {
	entities = make(map[string]*Entity)
}

// Add adds a entity to the map
func Add(e *Entity) {
	entities[e.String()] = e
}

// Get returns a entity from the map
func Get(id string) *Entity {
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
		result = append(result, e.GetAlias())
	}
	return result
}

func ListDIDs() []string {
	var result []string
	for _, e := range entities {
		result = append(result, e.String())
	}
	return result
}
