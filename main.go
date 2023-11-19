package main

import (
	"context"
	"fmt"
	"os"

	"github.com/bahner/go-home/actor"
	"github.com/bahner/go-home/config"
	"github.com/bahner/go-home/room"
	"github.com/bahner/go-ma/p2p"
	"github.com/libp2p/go-libp2p"

	log "github.com/sirupsen/logrus"
)

func main() {
	config.Init()

	ctx := context.Background()
	ctxTimeout, cancel := context.WithTimeout(ctx, config.GetDiscoveryTimeout())
	defer cancel()

	actorKeyset := config.GetActorKeyset()
	roomKeyset := config.GetRoomKeyset()
	cborData, _ := actorKeyset.IPNSKey.MarshalCBOR()
	fmt.Printf("actorKeyset: %s\n", cborData)
	os.Exit(0)

	log.Infof("Intializing actor with identity: %s", actorKeyset.IPNSKey.DID)

	// Conifgure libp2p from here only
	libp2pOpts := []libp2p.Option{
		libp2p.Identity(actorKeyset.IPNSKey.PrivKey),
	}

	node, ps, err := p2p.Init(ctxTimeout, libp2pOpts...)
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
	log.Debugf("Room initialized: %s", r.Entity.DID.Fragment)

	r.Enter(ps, a)

	// Draw the UI.
	log.Debugf("Starting text UI")
	ui := NewChatUI(ctx, node, ps, r, a)
	if err := ui.Run(); err != nil {
		log.Errorf("error running text UI: %s", err)
	}
}
