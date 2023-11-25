package p2p

import (
	"context"

	"github.com/bahner/go-ma-actor/p2p/dht"
	"github.com/bahner/go-ma-actor/p2p/mdns"
	"github.com/libp2p/go-libp2p/core/host"
	log "github.com/sirupsen/logrus"
)

func StartPeerDiscovery(ctx context.Context, h host.Host) error {
	log.Debug("Starting peer discovery...")

	done := make(chan struct{}, 2) // Buffered channel to avoid blocking

	// Start DHT discovery in a new goroutine
	go func() {
		dhtINstance, err := dht.Init(ctx, h)
		if err != nil {
			log.Errorf("Failed to initialise DHT. Peer discovery unsuccessful. %v ", err)
			done <- struct{}{} // Signal completion
			return
		}
		dht.DiscoverPeers(ctx, dhtINstance, h)
		done <- struct{}{} // Signal completion
	}()

	// Start MDNS discovery in a new goroutine
	go func() {
		mdns.DiscoverPeers(ctx, h)
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
