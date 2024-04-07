package actor

import (
	"fmt"

	"github.com/bahner/go-ma/key/set"
)

func NewFromPackedKeyset(data string, cached bool) (*Actor, error) {

	keyset, err := set.Unpack(data)
	if err != nil {
		return nil, fmt.Errorf("entity: failed to unpack keyset: %s", err)
	}

	return New(keyset)

}
