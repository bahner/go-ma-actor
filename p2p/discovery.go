package p2p

import (
	"context"
	"fmt"
	"time"

	"github.com/bahner/go-ma"
	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/p2p/mdns"
	log "github.com/sirupsen/logrus"
)

// DiscoverPeers starts the peer discovery process.
// The context should be cancelled to stop the discovery process, and
// should probably a timeout context.
// Host is just a libp2p host.
//
// DHT is a Kademlia DHT instance.
// If nil, a new DHT instance will be created.
// You might want to pass a DHT instance in Server mode here, for long running processes.
func (p *P2P) DiscoverPeers() error {

	ctx, cancel := config.P2PDiscoveryContext()
	defer cancel()

	// Start MDNS discovery in a new goroutine
	go func() {
		m, err := mdns.New(p.DHT.Host(), ma.RENDEZVOUS)
		if err != nil {
			log.Errorf("Failed to start MDNS discovery: %s", err)
			return
		}
		m.DiscoverPeers(ctx)
	}()

	// Wait for a discovery process to complete
	err := p.DHT.DiscoverPeers(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialise DHT. Peer discovery unsuccessful: %w", err)
	}

	return nil
}

// DiscoveryLoop is a blocking function that will periodically
// call DiscoverPeers() until the context is cancelled.
// This shouldn't be cancelled in normal operation.
// Each iteration will have a timeout of its own.

func (p *P2P) DiscoveryLoop(ctx context.Context) {
	log.Infof("Starting discovery with retry interval %s", config.P2PDiscoveryRetryIntervalString())
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := p.DiscoverPeers() // This will block until discovery is complete or timeout
			if err != nil {
				log.Debugf("Discovery attempt failed: %s", err)
			}
			sleepTime := config.P2PDiscoveryRetryInterval()
			log.Debugf("Discovery sleeping for %s", sleepTime.String())
			time.Sleep(sleepTime)
		}
	}
}
