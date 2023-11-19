package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/bahner/go-space/actor"

	"github.com/bahner/go-space/p2p/host"
	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	log "github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()
	var err error

	initConfig()
	log.Infof("Intializing actor with identity: %s", identity.IPNSKey.DID)

	// Create the node from the keyset.
	log.Debug("Creating p2p host from identity ...")
	node, err := host.New(
		libp2p.Identity(identity.IPNSKey.PrivKey),
		libp2p.ListenAddrStrings(),
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to create host: %v", err))
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

	a, err := actor.NewFromKeyset(identity, *forcePublish)
	if err != nil {
		panic(fmt.Sprintf("Failed to create actor: %v", err))
	}
	log.Infof("Actor initialized: %s", a.Entity.DID.Fragment)

	// Publish the identity to IPFS.

	r, err := NewRoom(room)
	if err != nil {
		panic(fmt.Sprintf("Failed to create room: %v", err))
	}

	r.Enter(a)

	// Draw the UI.
	ui := NewChatUI(ctx, r, a)
	if err := ui.Run(); err != nil {
		log.Errorf("error running text UI: %s", err)
	}
}
