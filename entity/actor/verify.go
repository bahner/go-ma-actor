package actor

import (
	"fmt"

	"github.com/bahner/go-ma-actor/entity"
)

func (a *Actor) IsValid() bool {

	return a.Verify() == nil

}

func (a *Actor) Verify() error {

	if a.Entity == nil {
		return entity.ErrEntityIsNil
	}

	err := a.Keyset.Verify()
	if err != nil {
		return fmt.Errorf("actor/veirfy: %w", err)
	}

	return a.Entity.Verify()

}
