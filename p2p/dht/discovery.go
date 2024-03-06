package dht

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma"
	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/p2p/peer"
	"github.com/libp2p/go-libp2p/core/discovery"
	p2peer "github.com/libp2p/go-libp2p/core/peer"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
	log "github.com/sirupsen/logrus"
)

var ErrAddrInfoAddrsEmpty = fmt.Errorf("addrinfo has no addresses")

// Takes a context and a DHT instance and discovers peers using the DHT.
// You might want to se server option or not for the DHT.
// Takes a variadic list of discovery.Option. You'll need this if you want to set a custom routing table.
func (d *DHT) DiscoverPeers(ctx context.Context, discoveryOpts ...discovery.Option) error {
	log.Debugf("Starting DHT peer discovery searching for peers with rendezvous string: %s", ma.RENDEZVOUS)

	log.Debugf("Number of open connections: %d", len(d.h.Network().Conns()))
	//  Trim connections
	log.Debugf("Trimming open connections to %d", config.P2PConnmgrLowWatermark())
	d.h.ConnManager().TrimOpenConns(ctx)

	log.Debugf("Peer discovery timeout: %v", config.P2PDiscoveryTimeout())
	log.Debugf("Peer discovery context %v", ctx)

	routingDiscovery := drouting.NewRoutingDiscovery(d.IpfsDHT)
	if routingDiscovery == nil {
		return fmt.Errorf("dht:discovery: failed to create routing discovery")
	}

	dutil.Advertise(ctx, routingDiscovery, ma.RENDEZVOUS, discoveryOpts...)
	log.Debugf("Advertising rendezvous string: %s", ma.RENDEZVOUS)

	peerChan, err := routingDiscovery.FindPeers(ctx, ma.RENDEZVOUS, discoveryOpts...)
	if err != nil {
		return fmt.Errorf("dht:discovery: peer discovery error: %w", err)
	}

	for p := range peerChan {
		if p.ID == d.h.ID() {
			continue // Skip self connection
		}

		// Check if the context was cancelled or timed out
		if ctx.Err() != nil {
			log.Warn("Context cancelled or timed out, stopping DHT peer discovery.")
			return ctx.Err()
		}

		// err := d.handleDiscoveredPeer(ctx, p)
		// if err != nil {
		// 	log.Warnf("Failed to protect discovered peer: %s: %v", p.ID.String(), err)
		// }
		go d.handleDiscoveredPeer(ctx, p)

	}

	return nil
}

// Protects a discovered peer by connecting to it and protecting it.
// Make sure to only called this function after a peer has been discovered
// and that it's allowed to connect to it.
func (d *DHT) handleDiscoveredPeer(ctx context.Context, pai p2peer.AddrInfo) error {
	log.Debugf("Discovered peer: %s", pai.ID.String())

	id := pai.ID.String()

	if !peer.IsAllowed(id) {
		// For this to happen the peer must be known AND have been disallowed
		// Or it is unknown and allowAll is false
		log.Debugf("Discovered peer %s is not allowed to connect", id)
		return nil
	}

	// Algo. If we get Addrs from the discover process, we use them and try to connect.
	// If the connection succeeds, we protect the peer and write the Addrs back to the peerstore.
	// If the connection fails, we try to fetch the Addrs from the peerstore and connect with them.

	if len(pai.Addrs) > 0 {
		log.Infof("Peer %s discovered with addresses, attempting to connect", id)
		p := peer.New(&pai)
		err := d.peerConnectAndUpdateIfSuccessful(ctx, p)
		if err == nil {
			log.Infof("Successfully discovered peer %s", id)
			return nil
		}
	}

	log.Warnf("Discovered peer %s has no addresses.", id)
	// Avoid creating any new peer objects until we have Addrs
	// Fetch it from the backend, if it exists.
	p, err := peer.Get(id)
	if err == nil {
		err = d.peerConnectAndUpdateIfSuccessful(ctx, p)
		if err == nil {
			log.Infof("Successfully discovered peer %s", id)
			return nil
		}
	}

	// When all else fails attempt to fetch the Addrs from the DHT
	log.Warnf("Failed to connect to peer %s, fetching addresses from DHT", id)
	a, err := d.FindPeer(ctx, pai.ID)
	if err != nil {
		log.Errorf("Failed to find peer %s: %v", id, err)
		return err
	}

	err = d.peerConnectAndUpdateIfSuccessful(ctx, peer.New(&a))
	if err != nil {
		log.Errorf("Failed to connect to peer %s: %v", id, err)
	}

	return err
}

func (d *DHT) peerConnectAndUpdateIfSuccessful(ctx context.Context, p peer.Peer) error {

	if len(p.AddrInfo.Addrs) == 0 {
		return ErrAddrInfoAddrsEmpty
	}

	err := d.h.Connect(ctx, *p.AddrInfo)
	if err != nil {
		log.Warnf("I would've unprotected: %s", p.ID)
		// d.h.ConnManager().UntagPeer(p.ID, ma.RENDEZVOUS)
		// d.h.ConnManager().Unprotect(p.ID, ma.RENDEZVOUS)
		return err
	}

	d.h.ConnManager().TagPeer(p.AddrInfo.ID, ma.RENDEZVOUS, defaultTagValue)
	d.h.ConnManager().Protect(p.AddrInfo.ID, ma.RENDEZVOUS)
	return peer.Set(p)

}
