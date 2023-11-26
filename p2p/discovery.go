package p2p

import (
	"context"

	"github.com/bahner/go-ma-actor/p2p/dht"
	"github.com/bahner/go-ma-actor/p2p/mdns"
	p2pdht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	log "github.com/sirupsen/logrus"
)

// StartPeerDiscovery starts the peer discovery process.
// The context should be cancelled to stop the discovery process, and
// should probably a timeout context.
// Host is just a libp2p host.
//
// DHT is a Kademlia DHT instance.
// If nil, a new DHT instance will be created.
// You might want to pass a DHT instance in Server mode here, for long running processes.
func StartPeerDiscovery(ctx context.Context, h host.Host, dhtInstance *p2pdht.IpfsDHT) error {
	log.Debug("Starting peer discovery...")
	var err error
	done := make(chan struct{}, 2) // Buffered channel to avoid blocking

	if dhtInstance == nil {
		dhtInstance, err = dht.Init(ctx, h)
		if err != nil {
			log.Errorf("Failed to initialise DHT. Peer discovery unsuccessful. %v ", err)
			done <- struct{}{} // Signal completion
			return err
		}
	}

	// Start DHT discovery in a new goroutine
	go func() {

		if err != nil {
			log.Errorf("Failed to initialise DHT. Peer discovery unsuccessful. %v ", err)
			done <- struct{}{} // Signal completion
			return
		}
		dht.DiscoverPeers(ctx, dhtInstance, h)
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
