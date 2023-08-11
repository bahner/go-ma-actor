package main

import (
	"context"
	"sync"

	"github.com/bahner/go-myspace/p2p/host"
	"github.com/bahner/go-myspace/p2p/pubsub"
)

func initPubSubService(ctx context.Context, wg *sync.WaitGroup, host *host.P2pHost) {

	defer wg.Done()

	// Start libp2p node and discover peers
	host.Init(ctx)
	host.StartPeerDiscovery(ctx)

	ps = pubsub.New(host)
	ps.Start(ctx)

}

func GetPubSubService() *pubsub.Service {
	return ps
}
