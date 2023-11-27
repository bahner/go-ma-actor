package p2p

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/p2p/connmgr"
	"github.com/bahner/go-ma-actor/p2p/node"
	"github.com/bahner/go-ma-actor/p2p/pubsub"
	"github.com/bahner/go-ma/key/ipns"
	libp2p "github.com/libp2p/go-libp2p"
	p2ppubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
)

var (
	ctxDiscovery context.Context
	cancel       context.CancelFunc

	n  host.Host
	ps *p2ppubsub.PubSub
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
// Also takes a variadic list of libp2p options.
// Of it's nil, an empty list will be used.
//
// The function return the libp2p node and a PubSub Service

func Init(ctx context.Context, i *ipns.Key, p2pOpts ...libp2p.Option) (host.Host, *p2ppubsub.PubSub, error) {

	// Initiate libp2p options, if none are provided
	if p2pOpts == nil {
		p2pOpts = []libp2p.Option{}
	}

	// Add the connection manager to the options
	connMgr, err := connmgr.Init()
	if err != nil {
		return nil, nil, fmt.Errorf("p2p.Init: failed to create connection manager: %w", err)
	}
	p2pOpts = append(p2pOpts, libp2p.ConnectionManager(connMgr))

	// Create a new libp2p Host that listens on a random TCP port
	n, err = node.New(i, p2pOpts...)
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

	discoveryTimeout := config.GetDiscoveryTimeout()
	ctxDiscovery, cancel = context.WithTimeout(ctx, discoveryTimeout)
	defer cancel()

	err = StartPeerDiscovery(ctxDiscovery, n, nil)
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
