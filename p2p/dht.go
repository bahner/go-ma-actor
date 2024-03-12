package p2p

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma"
	"github.com/bahner/go-ma-actor/p2p/connmgr"
	"github.com/bahner/go-ma-actor/p2p/peer"
	p2pDHT "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	p2peer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	log "github.com/sirupsen/logrus"
)

type DHT struct {
	*p2pDHT.IpfsDHT
	Host            host.Host
	ConnectionGater *connmgr.ConnectionGater
}

// Initialise The Kademlia DHT and bootstrap it.
// The context is used to abort the process, but context.Background() probably works fine.
// If nil is passed, a background context will be used.
//
// The host is a libp2p host.
//
// Takes a variadic list of dht.Option. You'll need this if you want to set a custom routing table.
// or set the mode to server. None is fine for ordinary use.

func NewDHT(h host.Host, cg *connmgr.ConnectionGater, dhtOpts ...p2pDHT.Option) (*DHT, error) {

	var err error

	d := &DHT{
		Host:            h,
		ConnectionGater: cg,
	}

	d.IpfsDHT, err = p2pDHT.New(context.Background(), h, dhtOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kademlia DHT: %w", err)
	}

	d.Bootstrap(context.Background())

	return d, nil
}

var (
	ErrFailedToCreateRoutingDiscovery = fmt.Errorf("failed to create routing discovery")
)

// Run a continuous discovery loop to find new peers
// The ctx should probably be a background context
func (d *DHT) discoveryLoop(ctx context.Context) error {
	routingDiscovery := routing.NewRoutingDiscovery(d.IpfsDHT)
	if routingDiscovery == nil {
		return ErrFailedToCreateRoutingDiscovery
	}

	peerChan, err := routingDiscovery.FindPeers(ctx, ma.RENDEZVOUS)
	if err != nil {
		return fmt.Errorf("peer discovery error: %w", err)
	}

	go advertisementLoop(ctx, routingDiscovery) // Run advertisement continuously in the background
	go discover(ctx, peerChan, d)               // Run discovery continuously in the background

	return nil
}

func discover(ctx context.Context, peerChan <-chan p2peer.AddrInfo, d *DHT) error {
	for {
		select {
		case p, ok := <-peerChan:
			if !ok {
				log.Fatalf("DHT peer channel closed, but it was supposed to be running in the background.")
			}

			if p.ID == d.Host.ID() {
				continue // Skip self
			}

			if err := peer.ConnectAndProtect(context.Background(), d.Host, p); err != nil {
				log.Warnf("Failed to connect to discovered DHT peer: %s: %v", p.ID.String(), err)
			}
		case <-ctx.Done():
			log.Debug("DHT discovery loop cancelled")
			return nil
		}
	}
}
