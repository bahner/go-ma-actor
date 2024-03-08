package dht

import (
	"context"
	"fmt"

	"time"

	"github.com/bahner/go-ma"
	"github.com/bahner/go-ma-actor/config"
	"github.com/libp2p/go-libp2p/core/discovery"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	"github.com/libp2p/go-libp2p/p2p/discovery/util"
	log "github.com/sirupsen/logrus"
)

var (
	ErrFailedToCreateRoutingDiscovery = fmt.Errorf("failed to create routing discovery")
)

// Run a continuous discovery loop to find new peers
// The ctx should probably be a background context
func (d *DHT) DiscoveryLoop(ctx context.Context) error {
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

func discover(ctx context.Context, peerChan <-chan peer.AddrInfo, d *DHT) error {
	for {
		select {
		case p, ok := <-peerChan:
			if !ok {
				log.Fatalf("DHT peer channel closed, but it was supposed to be running in the background.")
			}

			if p.ID == d.h.ID() {
				continue // Skip self
			}

			if err := d.PeerConnectAndUpdateIfSuccessful(context.Background(), p); err != nil {
				log.Warnf("Failed to connect to discovered peer: %s: %v", p.ID.String(), err)
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func advertisementLoop(ctx context.Context, routingDiscovery *routing.RoutingDiscovery, discoveryOpts ...discovery.Option) {

	ticker := time.NewTicker(config.P2PDiscoveryAdvertiseInterval())
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			util.Advertise(ctx, routingDiscovery, ma.RENDEZVOUS, discoveryOpts...)
			log.Debugf("Advertising rendezvous string: %s", ma.RENDEZVOUS)
		case <-ctx.Done():
			return
		}
	}
}
