package actor

import (
	"fmt"
)

func (a *Actor) IsValid() bool {

	return a.Verify() == nil

}

func (a *Actor) Verify() error {

	if a.Entity == nil {
		return fmt.Errorf("entity: no entity")
	}

	if a.Keyset == nil {
		return fmt.Errorf("entity: no keyset")
	}

	return a.Entity.Verify()

}
