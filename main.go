package main

import (
	"context"

	"github.com/bahner/go-ma-actor/actor"
	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/ui"
	"github.com/libp2p/go-libp2p/core/host"

	log "github.com/sirupsen/logrus"
)

var (
	ctx    context.Context
	node   host.Host
	a      *actor.Actor
	entity string
)

func init() {

	ctx = context.Background()
	node = config.GetNode()
	a = actor.GetActor()
	entity = config.GetEntity()

}

func main() {

	// Draw the UI.
	log.Debugf("Starting text UI")
	ui := ui.NewChatUI(ctx, node, a, entity)
	if err := ui.Run(); err != nil {
		log.Errorf("error running text UI: %s", err)
	}
}
