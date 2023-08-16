package main

import (
	"context"
	"sync"

	"github.com/bahner/go-space/p2p/host"
	"github.com/bahner/go-space/p2p/pubsub"
)

func createAndInitPubSubService(ctx context.Context, h *host.Host) (*pubsub.Service, error) {

	// Start libp2p node and discover peers
	h.Init(ctx)

	discoveryWg := &sync.WaitGroup{}

	discoveryWg.Add(2)
	go host.DiscoverDHTPeers(ctx, discoveryWg, h.Node, rendezvous)
	go host.DiscoverMDNSPeers(ctx, discoveryWg, h.Node, rendezvous)
	discoveryWg.Wait()

	ps := pubsub.New(h)
	ps.Start(ctx)

	return ps, nil
}
