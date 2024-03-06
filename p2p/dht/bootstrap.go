package dht

import (
	"context"
	"fmt"
	"sync"

	p2pDHT "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/peer"
	log "github.com/sirupsen/logrus"
)

// You should only call this once. It bootstraps the DHT to the network.
func (d *DHT) Bootstrap(ctx context.Context) error {
	log.Info("Initialising Kademlia DHT.")

	err := d.IpfsDHT.Bootstrap(context.Background())
	if err != nil {
		return fmt.Errorf("failed to bootstrap Kademlia DHT: %w", err)
	}
	log.Debug("Kademlia DHT bootstrap setup.")

	var wg sync.WaitGroup

	// Attempt to connect to all bootstrap peers
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

	// Reset the connection gater to its original allow state
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
