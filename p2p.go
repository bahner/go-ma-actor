package main

import (
	"context"
	"sync"

	"github.com/bahner/go-space/p2p/host"
	"github.com/bahner/go-space/p2p/pubsub"
	log "github.com/sirupsen/logrus"
)

func doDiscovery(ctx context.Context, h *host.Host) error {

	// Start libp2p node and discover peers
	h.Init(ctx)

	discoveryWg := &sync.WaitGroup{}

	// Discover peers
	// No need to log, as the discovery functions do that.
	discoveryWg.Add(1) // Only 1 of the following needs to finish
	go host.DiscoverDHTPeers(ctx, discoveryWg, h.Node, rendezvous)
	go host.DiscoverMDNSPeers(ctx, discoveryWg, h.Node, rendezvous)
	discoveryWg.Wait()

	return nil
}

func initSubscriptionService(ctx context.Context, h *host.Host) *pubsub.Service {

	doDiscovery(ctx, h)

	// Subscribe to the topic
	log.Debugf("Creating PubSub service ...")
	ps = pubsub.New(h)
	log.Debug("Starting PubSub service ...")
	ps.Start(ctx)
	log.Info("Subscription service started.")

	return ps
}
