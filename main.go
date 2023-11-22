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
	err error

	a *actor.Actor
	e string
	n host.Host
)

func init() {

	ctx = context.Background()

	// Try and run this as a goroutine. Not sure if it will work.
	n, _, err = p2p.Init(
		ctx,
		config.GetKeyset().IPNSKey,
		config.GetDiscoveryTimeout())
}

func main() {
	a, err = actor.NewFromKeyset(config.GetKeyset(), config.GetForcePublish())
	if err != nil {
		log.Fatalf("failed to create actor: %v", err)
	}

	e = config.GetEntity()
	// Draw the UI.
	log.Debugf("Starting text UI")
	ui := ui.NewChatUI(ctx, n, a, e)
	if err := ui.Run(); err != nil {
		log.Errorf("error running text UI: %s", err)
	}
}
