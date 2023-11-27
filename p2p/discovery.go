package p2p

import (
	"fmt"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/p2p/mdns"
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

	err := p.DHT.DiscoverPeers()
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
