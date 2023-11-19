package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/bahner/go-home/actor"
	"github.com/bahner/go-home/room"

	"github.com/bahner/go-space/p2p/host"
	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	log "github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()
	var err error

	initConfig()
	log.Infof("Intializing actor with identity: %s", actorKeyset.IPNSKey.DID)

	// Create the node from the keyset.
	log.Debug("Creating p2p host from identity ...")
	node, err := host.New(
		libp2p.Identity(actorKeyset.IPNSKey.PrivKey),
		libp2p.ListenAddrStrings(),
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to create p2p host: %v", err))
	}
	log.Debugf("node: %v", node)
	// the discoveryProcess return nil, so no need to check.
	log.Debug("Initializing subscription service ...")
	discoveryWg := &sync.WaitGroup{}

	// Discover peers
	// No need to log, as the discovery functions do that.
	discoveryWg.Add(1) // Only 1 of the following needs to finish
	go node.StartPeerDiscovery(ctx, discoveryWg, rendezvous)
	discoveryWg.Wait()

	ps, err = pubsub.NewGossipSub(ctx, node)
	if err != nil {
		panic(fmt.Sprintf("Failed to create pubsub service: %v", err))
	}

	a, err := actor.NewFromKeyset(ctx, ps, actorKeyset, *forcePublish)
	if err != nil {
		panic(fmt.Sprintf("Failed to create actor: %v", err))
	}
	log.Infof("Actor initialized: %s", a.Entity.DID.Fragment)

	ra, err := actor.NewFromKeyset(ctx, ps, roomKeyset, *forcePublish)
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
