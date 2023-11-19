package main

import (
	"context"
	"sync"

	"github.com/bahner/go-space/p2p/host"
	"github.com/bahner/go-space/p2p/pubsub"
)

func initSubscriptionService(ctx context.Context, h *host.P2pHost) *pubsub.Service {

	discoveryWg := &sync.WaitGroup{}

	// Discover peers
	// No need to log, as the discovery functions do that.
	discoveryWg.Add(1) // Only 1 of the following needs to finish
	go h.StartPeerDiscovery(ctx, discoveryWg, rendezvous)
	discoveryWg.Wait()

	return pubsub.New(ctx, h)
}
