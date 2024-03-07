package main

import (
	"fmt"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/p2p/connmgr"
	"github.com/bahner/go-ma-actor/p2p/dht"
	"github.com/bahner/go-ma-actor/p2p/node"
	"github.com/libp2p/go-libp2p"
)

func DHT(cg *connmgr.ConnectionGater) (*dht.DHT, error) {

	// THese are the relay specific parts.
	p2pOpts := []libp2p.Option{
		libp2p.ConnectionGater(cg),
	}

	n, err := node.New(config.NodeIdentity(), p2pOpts...)
	if err != nil {
		return nil, fmt.Errorf("pong: failed to create libp2p node: %w", err)
	}

	d, err := dht.New(n, cg)
	if err != nil {
		return nil, fmt.Errorf("pong: failed to create DHT: %w", err)
	}

	return d, nil
}
