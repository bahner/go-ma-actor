package ui

import (
	"fmt"

	"github.com/bahner/go-ma-actor/actor"
)

// Creates a UI, but panics if it fails.
func Init(a *actor.Actor) *ChatUI {
	fmt.Println("Creating text UI...")
	ui, err := New(a)
	if err != nil {
		panic(fmt.Sprintf("error creating text UI: %s", err))
	}
	return ui
}
