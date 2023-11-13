package main

import (
	"context"
	"fmt"

	"github.com/bahner/go-space/p2p/host"
	"github.com/libp2p/go-libp2p"
	log "github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()

	initConfig()
	log.Infof("Intializing actor with identity: %s", identity.IPNSKey.DID)

	// Create the node from the keyset.
	log.Debug("Creating p2p host from identity ...")
	node := host.New()
	node.AddOption(libp2p.Identity(identity.IPNSKey.PrivKey))
	// node.AddOption(libp2p.ListenAddrStrings(
	// 	"/ip4/0.0.0.0/tcp/0",
	// 	"/ip4/0.0.0.0/udp/0",
	// 	"/ip6/::/tcp/0",
	// 	"/ip6/::/udp/0"))
	log.Debugf("node: %v", node)
	// the discoveryProcess return nil, so no need to check.
	log.Debug("Initializing subscription service ...")
	ps = initSubscriptionService(ctx, node)

	a, err := initActor(identity)
	if err != nil {
		panic(fmt.Sprintf("Failed to create actor: %v", err))
	}
	log.Infof("Actor initialized: %s", a.DID.Fragment)

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
