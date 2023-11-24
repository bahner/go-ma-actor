package p2p

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bahner/go-ma-actor/p2p/node"
	"github.com/bahner/go-ma-actor/p2p/pubsub"
	"github.com/bahner/go-ma/key/ipns"
	p2ppubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
)

var (
	err error

	ctxDiscovery context.Context
	cancel       context.CancelFunc

	n  host.Host
	ps *p2ppubsub.PubSub

	connectedPeers = make(map[string]struct{})
	peerMutex      sync.Mutex
)

// Initialise everything needed for p2p communication.
//
// If ctx is nil a background context will be used as basis for a timeout context.
// So nil is fine.
//
// i is the ipns key.
//
// discoveryTimeout is the timeout duration for peer discovery.
// It's a time.Duration type.
//
// The function return the libp2p node and a PubSub Service

func Init(ctx context.Context, i *ipns.Key, discoveryTimeout time.Duration) (host.Host, *p2ppubsub.PubSub, error) {

	// Create a new libp2p Host that listens on a random TCP port
	n, err = node.New(i, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("p2p.Init: failed to create libp2p node: %w", err)
	}

	ps, err = pubsub.New(ctx, n)
	if err != nil {
		return nil, nil, fmt.Errorf("p2p.Init: failed to create pubsub: %w", err)
	}

	// Peer discovery
	if ctx == nil {
		ctx = context.Background()
	}

	ctxDiscovery, _ = context.WithTimeout(ctx, discoveryTimeout)
	// defer cancel()

	err = StartPeerDiscovery(ctxDiscovery, n)
	if err != nil {
		return nil, nil, fmt.Errorf("p2p.Init: failed to start peer discovery: %w", err)
	}

	return n, ps, nil
}

func GetPubSub() *p2ppubsub.PubSub {
	return ps
}

func GetNode() host.Host {
	return n
}

func GetConnectedPeers() []string {
	peerMutex.Lock()
	defer peerMutex.Unlock()

	peers := make([]string, 0, len(connectedPeers))
	for peer := range connectedPeers {
		peers = append(peers, peer)
	}
	return peers
}
