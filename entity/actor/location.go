package actor

import (
	"fmt"
)

func (e *Actor) GetLastKnownLocation() string {
	return e.Entity.Doc.LastKnownLocation
}

func (e *Actor) UpdateLastKnowLocation(location string) error {

	am, err := e.Entity.Doc.GetAssertionMethod()
	if err != nil {
		return fmt.Errorf("error getting assertion method: %w", err)
	}

	err = e.Entity.Doc.SetLastKnowLocation(location)
	if err != nil {
		return fmt.Errorf("error setting last known location: %w", err)
	}

	// Publish our new location
	e.Entity.Doc.UpdateVersion()
	err = e.Entity.Doc.Sign(e.Keyset.SigningKey, am)
	if err != nil {
		return fmt.Errorf("error signing document: %w", err)
	}

	// Spawn the publish in a goroutine so we don't block
	// It can take a while to publish
	go e.Entity.Doc.Publish()

	return nil
}
