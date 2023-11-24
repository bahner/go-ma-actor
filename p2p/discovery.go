package p2p

import (
	"context"

	"github.com/libp2p/go-libp2p/core/host"
	log "github.com/sirupsen/logrus"
)

func StartPeerDiscovery(ctx context.Context, h host.Host) error {
	log.Debug("Starting peer discovery...")

	done := make(chan struct{}, 2) // Buffered channel to avoid blocking

	// Start DHT discovery in a new goroutine
	go func() {
		DiscoverDHTPeers(ctx, h)
		done <- struct{}{} // Signal completion
	}()

	// Start MDNS discovery in a new goroutine
	go func() {
		DiscoverMDNSPeers(ctx, h)
		done <- struct{}{} // Signal completion
	}()

	// Wait for a discovery process to complete
	select {
	case <-ctx.Done():
		log.Warn("Peer discovery unsuccessful.")
		return ctx.Err()
	case <-done:
		log.Info("Peer discovery successful.")
		// Continue without waiting for the other process
		return nil
	}
}
