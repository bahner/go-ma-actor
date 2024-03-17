package main

import (
	"fmt"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/p2p"
	"github.com/bahner/go-ma-actor/p2p/connmgr"
	"github.com/bahner/go-ma-actor/p2p/node"
	"github.com/libp2p/go-libp2p"
	p2pDHT "github.com/libp2p/go-libp2p-kad-dht"
)

func DHT(cg *connmgr.ConnectionGater) (*p2p.DHT, error) {

	// THese are the relay specific parts.
	p2pOpts := []libp2p.Option{
		libp2p.ConnectionGater(cg),
		libp2p.Ping(true),
	}

	dhtOpts := []p2pDHT.Option{
		p2pDHT.Mode(p2pDHT.ModeServer),
	}

	n, err := node.New(config.NodeIdentity(), p2pOpts...)
	if err != nil {
		return nil, fmt.Errorf("pong: failed to create libp2p node: %w", err)
	}

	d, err := p2p.NewDHT(n, cg, dhtOpts...)
	if err != nil {
		return nil, fmt.Errorf("pong: failed to create DHT: %w", err)
	}

	return d, nil
}
