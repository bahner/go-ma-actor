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
