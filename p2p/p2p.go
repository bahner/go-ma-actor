package p2p

import (
	"context"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/p2p/node"
	"github.com/bahner/go-ma-actor/p2p/pubsub"
	p2ppubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	log "github.com/sirupsen/logrus"
)

var (
	err error

	ctxDiscovery context.Context
	cancel       context.CancelFunc

	n  host.Host
	ps *p2ppubsub.PubSub
)

func init() {

	// Set context timeout for peer discovery
	ctx := context.Background()

	discoveryTimeout := config.GetDiscoveryTimeout()
	ctxDiscovery, cancel = context.WithTimeout(ctx, discoveryTimeout)
	defer cancel()

	// Create a new libp2p Host that listens on a random TCP port
	n = node.Get()
	if err != nil {
		log.Fatalf("p2p: failed to create libp2p node: %v", err)
	}

	ps = pubsub.Get()

	// Start peer discovery
	err = StartPeerDiscovery(ctxDiscovery, n)
	if err != nil {
		log.Fatalf("p2p: failed to start peer discovery: %v", err)
	}

}

func GetPubSub() *p2ppubsub.PubSub {
	return ps
}

func GetNode() host.Host {
	return n
}
