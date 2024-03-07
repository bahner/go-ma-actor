package dht

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/bahner/go-ma"
	"github.com/bahner/go-ma-actor/config"
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
	// d.ConnectionGater.AllowAll = config.P2PDiscoveryAllow()

	for p := range peerChan {
		sem <- struct{}{} // Acquire a token
		go func(peerInfo p2peer.AddrInfo) {
			defer func() { <-sem }() // Release the token

			if peerInfo.ID == d.h.ID() {
				return // Skip self connection
			}
			if err := d.PeerConnectAndUpdateIfSuccessful(ctx, peerInfo); err != nil {
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
