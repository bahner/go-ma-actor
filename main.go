package main

import (
	"context"
	"fmt"

	"github.com/bahner/go-home/actor"
	"github.com/bahner/go-home/config"
	"github.com/bahner/go-home/room"

	log "github.com/sirupsen/logrus"
)

func main() {
	config.Init()

	ctx := context.Background()

	actorKeyset := config.GetActorKeyset()
	roomKeyset := config.GetRoomKeyset()

	log.Infof("Intializing actor with identity: %s", actorKeyset.IPNSKey.DID)

	ps, err := initPubSub(ctx, actorKeyset)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize pubsub: %v", err))
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

	r := room.Room{Actor: ra}

	r.Enter(ps, a)

	// // Draw the UI.
	// ui := NewChatUI(ctx, r, a)
	// if err := ui.Run(); err != nil {
	// 	log.Errorf("error running text UI: %s", err)
	// }
}
