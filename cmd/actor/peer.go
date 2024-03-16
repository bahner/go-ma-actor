package main

import (
	"fmt"

	"github.com/bahner/go-ma-actor/p2p/peer"
)

func initPeer(id string) error {

	p, err := peer.GetOrCreate(id)
	if err != nil {
		return fmt.Errorf("error getting or creating peer: %s", err)
	}

	return p.SetAllowed(true) // Ensure allowed even if we were previously denied.
}
