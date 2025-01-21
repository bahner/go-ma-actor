package actor

import (
	"fmt"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma-actor/p2p"
)

// Initialise an actor or panic, meant to be called from most main's.
// Panics if the actor is not valid.
// The parameters are P2P options, which should include the identity.
func Init(opts *p2p.Options) *Actor {
	fmt.Println("Creating actor from keyset...")
	a, err := New(config.ActorKeyset())
	if err != nil {
		panic(fmt.Sprintf("error creating actor: %s", err))
	}

	// Better safe than sorry.
	// Without a valid actor, we can't do anything.
	if a.Verify() != nil {
		panic(fmt.Sprintf("%s is not a valid actor: %v", a.Entity.DID.Id, err))
	}

	_, err = entity.Init(a.Entity.DID)
	if err != nil {
		panic(fmt.Sprintf("error getting or creating entity: %s", err))
	}

	// P2P start the p2p service
	fmt.Println("Initialising actor Host")
	a.P2P, err = p2p.Init(opts)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize p2p: %v", err))
	}

	return a
}
