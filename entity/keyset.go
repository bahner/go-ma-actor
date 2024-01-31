package entity

import (
	"fmt"

	"github.com/bahner/go-ma/key/set"
)

// Takes a keyset and an alias (name) and creates a new entity.
// The keyset is used to create the encryption and signing keys.
// The alias can be "" and will be set to the fragment of the DID.
func NewFromKeyset(k *set.Keyset, nick string) (*Entity, error) {

	return New(k.DID, k, nick)
}

func (e *Entity) IsValid() bool {

	return e.Verify() == nil

}

func NewFromPackedKeyset(data string) (*Entity, error) {

	keyset, err := set.Unpack(data)
	if err != nil {
		return nil, fmt.Errorf("entity: failed to unpack keyset: %s", err)
	}

	return NewFromKeyset(keyset, keyset.DID.Fragment)

}
