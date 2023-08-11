package main

import (
	"context"
	"sync"

	"github.com/bahner/go-myspace/p2p/host"
	"github.com/bahner/go-myspace/p2p/pubsub"
)

func initPubSubService(ctx context.Context, wg *sync.WaitGroup, h *host.P2pHost) {

	defer wg.Done()

	// Start libp2p node and discover peers
	h.Init(ctx)

	discoveryWg := &sync.WaitGroup{}

	discoveryWg.Add(2)
	go host.DiscoverDHTPeers(ctx, discoveryWg, h.Node, rendezvous)
	go host.DiscoverMDNSPeers(ctx, discoveryWg, h.Node, serviceName)
	discoveryWg.Wait()

	ps = pubsub.New(h)
	ps.Start(ctx)

}

func GetPubSubService() *pubsub.Service {
	return ps
}
