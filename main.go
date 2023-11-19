package main

import (
	"context"
	"fmt"

	"github.com/bahner/go-home/actor"
	"github.com/bahner/go-home/config"
	"github.com/bahner/go-home/p2p"
	"github.com/bahner/go-home/room"

	log "github.com/sirupsen/logrus"
)

func main() {
	config.Init()

	ctx := context.Background()

	actorKeyset := config.GetActorKeyset()
	roomKeyset := config.GetRoomKeyset()

	log.Infof("Intializing actor with identity: %s", actorKeyset.IPNSKey.DID)

	node, ps, err := p2p.Init(ctx, actorKeyset)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize p2p: %v", err))
	}

	a, err := actor.NewFromKeyset(ctx, ps, actorKeyset, config.GetForcePublish())
	if err != nil {
		panic(fmt.Sprintf("Failed to create actor: %v", err))
	}
	log.Infof("Actor initialized: %s", a.Entity.DID.Fragment)

	ra, err := actor.NewFromKeyset(ctx, ps, roomKeyset, config.GetForcePublish())
	if err != nil {
		panic(fmt.Sprintf("Failed to create room actor: %v", err))
	}

	r, err := room.New(ra)
	if err != nil {
		panic(fmt.Sprintf("Failed to create room: %v", err))
	}

	r.Enter(ps, a)

	// Draw the UI.
	ui := NewChatUI(ctx, node, ps, r, a)
	if err := ui.Run(); err != nil {
		log.Errorf("error running text UI: %s", err)
	}
}
