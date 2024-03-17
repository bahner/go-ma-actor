package main

import (
	"fmt"

	"github.com/bahner/go-ma-actor/entity/actor"
	"github.com/bahner/go-ma-actor/p2p"
	"github.com/bahner/go-ma-actor/ui"
)

func initUiOrPanic(p2P *p2p.P2P, a *actor.Actor) *ui.ChatUI {
	fmt.Println("Creating text UI...")
	ui, err := ui.New(p2P, a)
	if err != nil {
		panic(fmt.Sprintf("error creating text UI: %s", err))
	}
	return ui
}
