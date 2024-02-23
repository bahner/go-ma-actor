package dht

import (
	"context"
	"fmt"
	"sync"

	"github.com/bahner/go-ma"
	"github.com/bahner/go-ma-actor/config"
	p2pDHT "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/discovery"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
	log "github.com/sirupsen/logrus"
)

type DHT struct {
	*p2pDHT.IpfsDHT
	h host.Host
}

// Initialise The Kademlia DHT and bootstrap it.
// The context is used to abort the process, but context.Background() probably works fine.
// If nil is passed, a background context will be used.
//
// The host is a libp2p host.
//
// Takes a variadic list of dht.Option. You'll need this if you want to set a custom routing table.
// or set the mode to server. None is fine for ordinary use.

func New(h host.Host, dhtOpts ...p2pDHT.Option) (*DHT, error) {

	var err error

	d := &DHT{h: h}

	d.IpfsDHT, err = p2pDHT.New(context.Background(), h, dhtOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kademlia DHT: %w", err)
	}

	d.Bootstrap(context.Background())

	return d, nil
}

func (d *DHT) Bootstrap(ctx context.Context) error {
	log.Info("Initialising Kademlia DHT.")

	err := d.IpfsDHT.Bootstrap(context.Background())
	if err != nil {
		return fmt.Errorf("failed to bootstrap Kademlia DHT: %w", err)
	}
	log.Debug("Kademlia DHT bootstrap setup.")

	var wg sync.WaitGroup
	for _, peerAddr := range p2pDHT.DefaultBootstrapPeers {
		peerinfo, err := peer.AddrInfoFromP2pAddr(peerAddr)
		if err != nil {
			log.Warnf("Failed to convert bootstrap peer address: %v", err)
			continue
		}

		log.Debugf("Bootstrapping to peer: %s", peerinfo.ID.String())

		wg.Add(1)
		go func(pInfo peer.AddrInfo) {
			defer wg.Done()

			if err := d.h.Connect(ctx, pInfo); err != nil {
				log.Warnf("Bootstrap warning: %v", err)
			}
		}(*peerinfo)
	}

	// Wait for all bootstrap attempts to complete or context cancellation
	wg.Wait()
	log.Info("All bootstrap attempts completed.")

	select {
	case <-ctx.Done():
		log.Warn("Context cancelled during bootstrap.")
		return ctx.Err()
	default:
		log.Info("Kademlia DHT bootstrapped successfully.")
		return nil
	}
}

// Takes a context and a DHT instance and discovers peers using the DHT.
// You might want to se server option or not for the DHT.
// Takes a variadic list of discovery.Option. You'll need this if you want to set a custom routing table.
func (d *DHT) DiscoverPeers(ctx context.Context, discoveryOpts ...discovery.Option) error {
	log.Debugf("Starting DHT peer discovery searching for peers with rendezvous string: %s", ma.RENDEZVOUS)

	log.Debugf("Number of open connections: %d", len(d.h.Network().Conns()))
	//  Trim connections
	log.Debugf("Trimming open connections to %d", config.GetLowWatermark())
	d.h.ConnManager().TrimOpenConns(ctx)

	log.Debugf("Peer discovery timeout: %v", config.GetDiscoveryTimeout())
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

		err := d.h.Connect(ctx, p)
		if err != nil {
			log.Debugf("Failed connecting to %s, error: %v", p.ID.String(), err)
			// d.h.ConnManager().UntagPeer(p.ID, ma.RENDEZVOUS)
			// d.h.ConnManager().Unprotect(p.ID, ma.RENDEZVOUS)
			continue
		}

		log.Debugf("Connected to DHT peer: %s", p.ID.String())
		d.h.ConnManager().TagPeer(p.ID, ma.RENDEZVOUS, 10)
		d.h.ConnManager().Protect(p.ID, ma.RENDEZVOUS)

		// Check if the context was cancelled or timed out
		if ctx.Err() != nil {
			log.Warn("Context cancelled or timed out, stopping DHT peer discovery.")
			return ctx.Err()
		}

	}

	return nil
}
