package entity

import "fmt"

func (e *Entity) Verify() error {

	if e.Ctx == nil {
		return fmt.Errorf("entity/verify: context is nil")
	}

	if e.CancelFunc == nil {
		return fmt.Errorf("entity/verify: cancel is nil")
	}

	if e.DID == nil {
		return fmt.Errorf("entity/verify: did is nil")
	}

	if e.Doc == nil {
		return fmt.Errorf("entity/verify: document is nil")
	}

	if e.Topic == nil {
		return fmt.Errorf("entity/verify: topic is nil")
	}

	if !e.DID.IsValid() {
		return fmt.Errorf("entity/verify: did is invalid")
	}

	err := e.Doc.Verify()
	if err != nil {
		return fmt.Errorf("entity/verify: failed to verify document: %w", err)
	}

	return nil
}
