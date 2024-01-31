package entity

import (
	"fmt"

	"github.com/bahner/go-ma/did/doc"
)

func (e *Entity) GetLastKnownLocation() string {
	return e.Doc.LastKnownLocation
}

func (e *Entity) UpdateLastKnowLocation(location string) error {

	am, err := e.Doc.GetAssertionMethod()
	if err != nil {
		return fmt.Errorf("error getting assertion method: %w", err)
	}

	err = e.Doc.SetLastKnowLocation(location)
	if err != nil {
		return fmt.Errorf("error setting last known location: %w", err)
	}

	// Publish our new location
	e.Doc.UpdateVersion()
	err = e.Doc.Sign(e.Keyset.SigningKey, am)
	if err != nil {
		return fmt.Errorf("error signing document: %w", err)
	}

	// Spawn the publish in a goroutine so we don't block
	// It can take a while to publish
	opts := doc.DefaultPublishOptions()
	opts.Force = true
	go e.Doc.Publish(opts)

	return nil
}
