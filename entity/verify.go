package entity

import (
	"fmt"
)

func (e *Entity) Verify() error {

	err := e.DID.Verify()
	if err != nil {
		return fmt.Errorf("entity/verify: %w", err)
	}

	if e.Topic == nil {
		return fmt.Errorf("entity/verify: %w", ErrTopicIsNil)
	}

	err = e.Doc.Verify()
	if err != nil {
		return fmt.Errorf("entity/verify: %w", err)
	}

	return nil
}
