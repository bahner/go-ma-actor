package dht

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/p2p/connmgr"
	p2pDHT "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
)

type DHT struct {
	*p2pDHT.IpfsDHT
	h               host.Host
	ConnectionGater *connmgr.ConnectionGater
}

var defaultTagValue = 100

// Initialise The Kademlia DHT and bootstrap it.
// The context is used to abort the process, but context.Background() probably works fine.
// If nil is passed, a background context will be used.
//
// The host is a libp2p host.
//
// Takes a variadic list of dht.Option. You'll need this if you want to set a custom routing table.
// or set the mode to server. None is fine for ordinary use.

func New(h host.Host, cg *connmgr.ConnectionGater, dhtOpts ...p2pDHT.Option) (*DHT, error) {

	var err error

	d := &DHT{
		h:               h,
		ConnectionGater: cg,
	}

	d.IpfsDHT, err = p2pDHT.New(context.Background(), h, dhtOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kademlia DHT: %w", err)
	}

	d.ConnectionGater.AllowAll = true
	d.Bootstrap(context.Background())
	// Reset the connection gater to its original allow state
	d.ConnectionGater.AllowAll = config.P2PDiscoveryAllowAll()

	return d, nil
}
