package p2p

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma-actor/p2p/node"
	"github.com/bahner/go-ma-actor/p2p/pubsub"
	"github.com/libp2p/go-libp2p"
	p2ppubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
)

// This package initializes the p2p node and pubsub service, and starts peer discovery.
// That's should be all you need to get started with p2p in your application.
// The context is for peer discovery, and can be canceled to stop peer discovery.
// If set to nil, a background context will be used.
// The opts parameter is passed to the libp2p node constructor.
func Init(discoveryCtx context.Context, opts ...libp2p.Option) (host.Host, *p2ppubsub.PubSub, error) {

	if discoveryCtx == nil {
		discoveryCtx = context.Background()
	}

	// Create a new libp2p Host that listens on a random TCP port
	n, err := node.New(opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("p2p: failed to create libp2p node: %v", err)
	}

	// Start peer discovery
	err = StartPeerDiscovery(discoveryCtx, n)
	if err != nil {
		return nil, nil, fmt.Errorf("p2p: failed to start peer discovery: %v", err)
	}

	// Initialize pubsub service
	ps, err := pubsub.New(context.Background(), n)
	if err != nil {
		return nil, nil, fmt.Errorf("p2p: failed to create pubsub service: %v", err)
	}

	return n, ps, nil
}
