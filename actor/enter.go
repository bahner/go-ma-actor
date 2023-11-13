package actor

import "fmt"

// Takes a room topic and joins it. The room is the DID of the room actor.
func (a *Actor) Enter(room string) error {

	var err error

	// First close the current subscription
	a.Public.Close()

	a.Public, err = a.ps.Join(room)
	if err != nil {
		return fmt.Errorf("home: %v failed to join topic: %v", a, err)
	}

	return nil

}
