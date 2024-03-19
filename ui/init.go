package ui

import (
	"fmt"

	"github.com/bahner/go-ma-actor/entity/actor"
	"github.com/bahner/go-ma-actor/p2p"
)

func Init(p2P *p2p.P2P, a *actor.Actor) *ChatUI {
	fmt.Println("Creating text UI...")
	ui, err := New(p2P, a)
	if err != nil {
		panic(fmt.Sprintf("error creating text UI: %s", err))
	}
	return ui
}
