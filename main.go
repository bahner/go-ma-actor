package main

import (
	"context"

	"github.com/bahner/go-ma-actor/actor"
	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/p2p"
	"github.com/bahner/go-ma-actor/ui"
	"github.com/libp2p/go-libp2p/core/host"

	log "github.com/sirupsen/logrus"
)

var (
	ctx context.Context

	a *actor.Actor
	e string
	n host.Host
)

func init() {

	ctx = context.Background()

	a = actor.GetActor()
	e = config.GetEntity()
	n = p2p.GetNode()

}

func main() {

	// Draw the UI.
	log.Debugf("Starting text UI")
	ui := ui.NewChatUI(ctx, n, a, e)
	if err := ui.Run(); err != nil {
		log.Errorf("error running text UI: %s", err)
	}
}
