package pubsub

import (
	"context"
	"fmt"
	"sync"

	"github.com/bahner/go-home/config"
	"github.com/bahner/go-ma/key/set"
	"github.com/bahner/go-space/p2p/host"
	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	log "github.com/sirupsen/logrus"
)

func Init(ctx context.Context, k *set.Keyset) (*pubsub.PubSub, error) {

	// Create the node from the keyset.

	log.Debug("Creating p2p host from identity ...")
	node, err := host.New(
		libp2p.Identity(k.IPNSKey.PrivKey),
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
	go node.StartPeerDiscovery(ctx, discoveryWg, config.GetRendezvous())
	log.Debug("Waiting for discovery to finish ...")
	discoveryWg.Wait()

	return pubsub.NewGossipSub(ctx, node)
}
