package p2p

import (
	"context"
	"fmt"
	"time"

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

	ctx, cancel := config.GetDiscoveryContext()
	defer cancel()

	err := p.DHT.DiscoverPeers(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialise DHT. Peer discovery unsuccessful: %w", err)
	}

	// Start MDNS discovery in a new goroutine
	go func() {
		mdns.DiscoverPeers(ctx, p.Node)
	}()

	// Wait for a discovery process to complete

	return nil
}

// DiscoveryLoop is a blocking function that will periodically
// call DiscoverPeers() until the context is cancelled.
// This shouldn't be cancelled in normal operation.
// Each iteration will have a timeout of its own.

func (p *P2P) DiscoveryLoop(ctx context.Context) {
	log.Infof("Starting discovery with retry interval %s", config.GetDiscoveryRetryIntervalString())
	for {
		select {
		case <-ctx.Done():
			return
		default:
			p.DHT.DiscoverPeers(ctx)
			sleepTime := config.GetDiscoveryRetryInterval()
			log.Debugf("Sleeping for %s", sleepTime.String())
			time.Sleep(sleepTime)
		}
	}
}
