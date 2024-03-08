package p2p

import (
	"context"
)

// DiscoverPeers starts the peer discovery process.
// The context should be cancelled to stop the discovery process, and
// should probably a timeout context.
// Host is just a libp2p host.
//
// DHT is a Kademlia DHT instance.
// If nil, a new DHT instance will be created.
// You might want to pass a DHT instance in Server mode here, for long running processes.
func (p *P2P) DiscoveryLoop(ctx context.Context) error {

	go p.MDNS.DiscoveryLoop(ctx)

	go p.DHT.DiscoveryLoop(ctx)

	return nil
}
