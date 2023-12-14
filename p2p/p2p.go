package p2p

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/p2p/dht"
	"github.com/bahner/go-ma-actor/p2p/node"
	"github.com/bahner/go-ma-actor/p2p/pubsub"
	libp2p "github.com/libp2p/go-libp2p"
	p2ppubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
)

// This is not a normal libp2p node, it's a wrapper around it. And it is specific to this project.
// It contains a libp2p node, a pubsub service and a DHT instance.
// It also contains a list of connected peers.
type P2P struct {
	Node   host.Host
	PubSub *p2ppubsub.PubSub
	DHT    *dht.DHT
}

// Initialise everything needed for p2p communication. The function forces use of a specific IPNS key.
// Taken from the config package. It would be an error to initialise the node with a different key.
//
// d is a DHT instance. If nil, a new one will be created.
//
// Also takes a variadic list of libp2p options.
// Of it's nil, an empty list will be used.
// You want to add your connectionmanager to
//
// The function return the libp2p node and a PubSub Service

func Init(d *dht.DHT, p2pOpts ...libp2p.Option) (*P2P, error) {

	// Create a new libp2p Host that listens on a random TCP port
	n, err := node.New(config.GetNodeIdentity(), p2pOpts...)
	if err != nil {
		return nil, fmt.Errorf("p2p.Init: failed to create libp2p node: %w", err)
	}

	ps, err := pubsub.New(context.Background(), n)
	if err != nil {
		return nil, fmt.Errorf("p2p.Init: failed to create pubsub: %w", err)
	}

	if d == nil {

		d, err = dht.New(n)
		if err != nil {
			return nil, fmt.Errorf("p2p.Init: failed to create DHT: %w", err)
		}
	}

	myP2P := &P2P{
		Node:   n,
		PubSub: ps,
		DHT:    d,
	}

	return myP2P, nil
}
