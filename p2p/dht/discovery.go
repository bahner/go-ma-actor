package dht

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/bahner/go-ma"
	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/p2p/peer"
	"github.com/libp2p/go-libp2p/core/discovery"
	"github.com/libp2p/go-libp2p/core/network"
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
	// //  Trim connections
	// log.Debugf("Trimming open connections to %d", config.P2PConnmgrLowWatermark())
	// d.h.ConnManager().TrimOpenConns(ctx)

	log.Debugf("Peer discovery timeout: %v", config.P2PDiscoveryTimeout())
	log.Debugf("Peer discovery context %v", ctx)

	routingDiscovery := drouting.NewRoutingDiscovery(d.IpfsDHT)
	if routingDiscovery == nil {
		return fmt.Errorf("dht:discovery: failed to create routing discovery")
	}

	// discoveryOpts = append(discoveryOpts,
	// 	discovery.Limit(config.P2PDiscoveryAdvertiseLimit()),
	// 	discovery.TTL(config.P2PDiscoveryAdvertiseTTL()))

	dutil.Advertise(ctx, routingDiscovery, ma.RENDEZVOUS, discoveryOpts...)
	log.Debugf("Advertising rendezvous string: %s", ma.RENDEZVOUS)

	peerChan, err := routingDiscovery.FindPeers(ctx, ma.RENDEZVOUS, discoveryOpts...)
	if err != nil {
		return fmt.Errorf("dht:discovery: peer discovery error: %w", err)
	}

	sem := make(chan struct{}, config.P2PDiscoveryLimit()) // Semaphore for controlling concurrency
	var successCount int32

	// // Make sure we have set the allowAll flag to it's to it's allowed state
	// d.ConnectionGater.AllowAll = config.P2PDiscoveryAllowAll()

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

	// If peer is explicitly denied, close any connections
	if !peer.IsAllowed(id) {
		// For this to happen the peer must be known AND have been disallowed
		// Or it is unknown and allowAll is false
		log.Warnf("Discovered peer %s has been denied access.", id)
		log.Warnf("Closing connections to peer %s", id)
		d.h.Network().ClosePeer(pai.ID)
		return nil
	}

	// If the peer is already connected, skip
	if d.h.Network().Connectedness(pai.ID) == network.Connected {
		log.Debugf("Peer %s is already connected", id)
		return nil
	}

	p, err := peer.GetOrCreateFromAddrInfo(&pai)
	if err != nil {
		return err
	}

	if len(pai.Addrs) > 0 {
		log.Debugf("Peer %s discovered with addresses, attempting to connect", id)
		err = d.PeerConnectAndUpdateIfSuccessful(ctx, p)
		if err != nil {
			log.Debugf("Failed to connect to newly discovered peer %s: %v", id, err)
			return err
		}
	}

	log.Infof("Successfully discovered peer %s", id)
	return nil
}

func (d *DHT) PeerConnectAndUpdateIfSuccessful(ctx context.Context, p peer.Peer) error {

	if len(p.AddrInfo.Addrs) == 0 {
		return ErrAddrInfoAddrsEmpty
	}

	id := p.AddrInfo.ID

	err := d.h.Connect(ctx, *p.AddrInfo)
	// NOOP. Clients that are protected are allowed to connect to us.
	// Even if we can't connect to them right now, we should still protect them.
	// if err != nil && d.h.ConnManager().IsProtected(id, ma.RENDEZVOUS) {
	// log.Warnf("Unprotecting previously protected peer %s: %v", id, err)
	// d.h.ConnManager().UntagPeer(id, ma.RENDEZVOUS)
	// d.h.ConnManager().Unprotect(id, ma.RENDEZVOUS)
	// }
	if err != nil {
		return err
	}

	if !d.h.ConnManager().IsProtected(id, ma.RENDEZVOUS) {
		log.Infof("Protecting previously unprotected peer %s", id)
		d.h.ConnManager().TagPeer(p.AddrInfo.ID, ma.RENDEZVOUS, defaultTagValue)
		d.h.ConnManager().Protect(p.AddrInfo.ID, ma.RENDEZVOUS)

		// This is a new peer, so we should allow it explicitly.
		// ACtually it should be allowed by default, but we'll set it explicitly here.
		// Ref. line #99 above
		p.Allowed = true

	}

	return peer.Set(p)

}
