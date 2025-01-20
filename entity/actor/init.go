package actor

import (
	"fmt"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma-actor/p2p"
)

// Initialise an actor or panic.
// This is a common sugar function to create an actor from the keyset and set the DID Document.
// Meant to be called from most main's.
// Panics if the actor is not valid.
func Init(opts *p2p.Options) *Actor {
	fmt.Println("Creating actor from keyset...")
	a, err := New(config.ActorKeyset())
	if err != nil {
		panic(fmt.Sprintf("error creating actor: %s", err))
	}

	if opts == nil {
		opts = p2p.DefaultOptions()
	}

	id := a.Entity.DID.Id

	// Better safe than sorry.
	// Without a valid actor, we can't do anything.
	if a.Verify() != nil {
		panic(fmt.Sprintf("%s is not a valid actor: %v", id, err))
	}

	_, err = entity.Fetch(a.Entity.DID)
	if err != nil {
		panic(fmt.Sprintf("error getting or creating entity: %s", err))
	}

	// P2P start the p2p service

	return a
}
