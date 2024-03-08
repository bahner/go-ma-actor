package dht

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	log "github.com/sirupsen/logrus"
)

var (
	ErrFailedToCreateRoutingDiscovery = fmt.Errorf("failed to create routing discovery")
	ErrPeerChanClosed                 = fmt.Errorf("peer channel closed")
)

// Run a continuous discovery loop to find new peers
// The ctx should probably be a background context
func (d *DHT) DiscoveryLoop(ctx context.Context) error {
	routingDiscovery := drouting.NewRoutingDiscovery(d.IpfsDHT)
	if routingDiscovery == nil {
		return ErrFailedToCreateRoutingDiscovery
	}

	peerChan, err := routingDiscovery.FindPeers(ctx, ma.RENDEZVOUS)
	if err != nil {
		return fmt.Errorf("peer discovery error: %w", err)
	}

	for {
		select {
		case p, ok := <-peerChan:
			if !ok {
				if !(ctx == nil) {
					log.Fatalf("DHT peer channel closed, but it was supposed to be running in the background.")
				}
				return ErrPeerChanClosed
			}

			if p.ID == d.h.ID() {
				continue // Skip self
			}

			if err := d.PeerConnectAndUpdateIfSuccessful(ctx, p); err != nil {
				log.Warnf("Failed to connect to discovered peer: %s: %v", p.ID.String(), err)
			}
		case <-ctx.Done():
			return nil
		}
	}
}
