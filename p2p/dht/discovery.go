package dht

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/bahner/go-ma"
	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/p2p/peer"
	"github.com/libp2p/go-libp2p/core/discovery"
	p2peer "github.com/libp2p/go-libp2p/core/peer"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
	log "github.com/sirupsen/logrus"
)

var (
	ErrAddrInfoAddrsEmpty    = fmt.Errorf("addrinfo has no addresses")
	ErrNoProtectedPeersFound = fmt.Errorf("no peers were discovered")
)

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

	discoveryOpts = append(discoveryOpts,
		discovery.Limit(config.P2PDiscoveryAdvertiseLimit()),
		discovery.TTL(config.P2PDiscoveryAdvertiseTTL()))

	dutil.Advertise(ctx, routingDiscovery, ma.RENDEZVOUS, discoveryOpts...)
	log.Debugf("Advertising rendezvous string: %s", ma.RENDEZVOUS)

	peerChan, err := routingDiscovery.FindPeers(ctx, ma.RENDEZVOUS, discoveryOpts...)
	if err != nil {
		return fmt.Errorf("dht:discovery: peer discovery error: %w", err)
	}

	sem := make(chan struct{}, config.P2PDiscoveryLimit()) // Semaphore for controlling concurrency
	var successCount int32                                 // Atomic counter for tracking successful connections

	for p := range peerChan {
		sem <- struct{}{} // Acquire a token
		go func(peerInfo p2peer.AddrInfo) {
			defer func() { <-sem }() // Release the token

			if peerInfo.ID == d.h.ID() {
				return // Skip self connection
			}
			if err := d.handleDiscoveredPeer(ctx, peerInfo); err != nil {
				log.Warnf("Failed to protect discovered peer: %s: %v", peerInfo.ID.String(), err)
			} else {
				atomic.AddInt32(&successCount, 1) // Increment on successful operation
			}
		}(p)
	}

	// Wait for all goroutines to finish
	for i := 0; i < cap(sem); i++ {
		sem <- struct{}{} // Ensure all tokens are returned before proceeding
	}

	// After processing all peers, check if there were any successful connections
	if atomic.LoadInt32(&successCount) == 0 {
		return ErrNoProtectedPeersFound
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
		log.Warnf("Discovered peer %s has been denied access.", id)
		return nil
	}

	// Algo. If we get Addrs from the discover process, we use them and try to connect.
	// If the connection succeeds, we protect the peer and write the Addrs back to the peerstore.
	// If the connection fails, we try to fetch the Addrs from the peerstore and connect with them.

	if len(pai.Addrs) > 0 {
		log.Debugf("Peer %s discovered with addresses, attempting to connect", id)
		p := peer.New(&pai)
		err := d.peerConnectAndUpdateIfSuccessful(ctx, p)
		if err == nil {
			log.Infof("Successfully discovered peer %s", id)
			return nil
		}
	}

	log.Debugf("Discovered peer %s has no addresses.", id)
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
	log.Debugf("Failed to connect to peer %s, fetching addresses from DHT", id)
	a, err := d.FindPeer(ctx, pai.ID)
	if err != nil {
		log.Debugf("Failed to find peer %s: %v", id, err)
		return err
	}

	err = d.peerConnectAndUpdateIfSuccessful(ctx, peer.New(&a))
	if err != nil {
		log.Debugf("Failed to connect to peer %s: %v", id, err)
	}

	return err
}

func (d *DHT) peerConnectAndUpdateIfSuccessful(ctx context.Context, p peer.Peer) error {

	if len(p.AddrInfo.Addrs) == 0 {
		return ErrAddrInfoAddrsEmpty
	}

	id := p.AddrInfo.ID

	err := d.h.Connect(ctx, *p.AddrInfo)
	if err != nil && d.h.ConnManager().IsProtected(id, ma.RENDEZVOUS) {
		log.Warnf("Unprotecting previously protected peer %s: %v", id, err)
		d.h.ConnManager().UntagPeer(id, ma.RENDEZVOUS)
		d.h.ConnManager().Unprotect(id, ma.RENDEZVOUS)
	}

	if err != nil {
		return err
	}
	if !d.h.ConnManager().IsProtected(id, ma.RENDEZVOUS) {
		log.Infof("Protecting previously unprotected peer %s", id)
		d.h.ConnManager().TagPeer(p.AddrInfo.ID, ma.RENDEZVOUS, defaultTagValue)
		d.h.ConnManager().Protect(p.AddrInfo.ID, ma.RENDEZVOUS)
	}
	return peer.Set(p)

}
