package space

import "fmt"

var spaces = map[string]*Space{}

// Fetch a space if it is known.
func Get(id string) *Space {
	return spaces[id]
}

// Delete a space if it is known.
// Will fetch it from the spaces map, close it and delete it.
func Delete(id string) error {

	s := spaces[id]
	if s == nil {
		return fmt.Errorf("space %s does not exist", id)
	}

	s.Close()

	delete(spaces, id)

	return nil
}

func CloseAll() {
	for _, s := range spaces {
		s.Close()
	}
}

func List() []string {
	var list []string
	for id := range spaces {
		list = append(list, id)
	}
	return list
}

// GetOrCreate a space if it is known.
// Will fetch it from the spaces map or create it.
func GetOrCreate(id string) (*Space, error) {

	s := spaces[id]
	if s != nil {
		return s, nil
	}

	return Create(id)
}

// Create a space if it is known.
// Will create it and add it to the spaces map.
func Create(id string) (*Space, error) {

	s, err := New(id)
	if err != nil {
		return nil, err
	}

	spaces[id] = s

	return s, nil
}
