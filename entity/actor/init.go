package actor

import (
	"fmt"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/entity"
)

// Initialise an actor or panic.
// This is a common sugar function to create an actor from the keyset and set the DID Document.
// Meant to be called from most main's.
// Panics if the actor is not valid.
func Init() *Actor {
	// The actor is needed for initialisation of the WebHandler.
	fmt.Println("Creating actor from keyset...")
	a, err := NewFromKeyset(config.ActorKeyset())
	if err != nil {
		panic(fmt.Sprintf("error creating actor: %s", err))
	}

	id := a.Entity.DID.Id

	fmt.Println("Creating and setting DID Document for actor...")
	err = a.CreateAndSetEntityDocument(id)
	if err != nil {
		panic(fmt.Sprintf("error creating document: %s", err))
	}

	// Better safe than sorry.
	// Without a valid actor, we can't do anything.
	if a == nil || a.Verify() != nil {
		panic(fmt.Sprintf("%s is not a valid actor: %v", id, err))
	}

	_, err = entity.New(a.Entity.DID)
	if err != nil {
		panic(fmt.Sprintf("error getting or creating entity: %s", err))
	}

	return a
}
