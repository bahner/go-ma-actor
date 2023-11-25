package p2p

import (
	"context"
	"fmt"
	"sync"

	"github.com/bahner/go-ma"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
	log "github.com/sirupsen/logrus"
)

func initDHT(ctx context.Context, h host.Host) (*dht.IpfsDHT, error) {
	log.Info("Initializing DHT.")

	kademliaDHT, err := dht.New(ctx, h)
	if err != nil {
		log.Error("Failed to create Kademlia DHT.")
		return nil, err
	} else {
		log.Debug("Kademlia DHT created.")
	}

	err = kademliaDHT.Bootstrap(ctx)
	if err != nil {
		log.Error("Failed to bootstrap Kademlia DHT.")
		return nil, err
	} else {
		log.Debug("Kademlia DHT bootstrap setup.")
	}

	var wg sync.WaitGroup
	for _, peerAddr := range dht.DefaultBootstrapPeers {
		peerinfo, err := peer.AddrInfoFromP2pAddr(peerAddr)
		if err != nil {
			log.Warnf("Failed to convert bootstrap peer address: %v", err)
			continue
		}

		log.Debugf("Bootstrapping to peer: %s", peerinfo.ID.String())

		wg.Add(1)
		go func(pInfo peer.AddrInfo) {
			defer wg.Done()

			errCh := make(chan error, 1)
			go func() {
				errCh <- h.Connect(ctx, pInfo)
			}()

			select {
			case <-ctx.Done():
				log.Debug("Context cancelled, aborting connection attempt.")
				return
			case err := <-errCh:
				if err != nil {
					log.Warnf("Bootstrap warning: %v", err)
				}
			}
		}(*peerinfo)
	}

	// Wait for all bootstrap attempts to complete or context cancellation
	go func() {
		wg.Wait()
		log.Info("All bootstrap attempts completed.")
	}()

	select {
	case <-ctx.Done():
		log.Info("Context cancelled during bootstrap.")
		return nil, ctx.Err()
	default:
		// Continue with other operations if context is not cancelled
	}

	log.Info("Kademlia DHT bootstrapped successfully.")
	return kademliaDHT, nil
}
func DiscoverDHTPeers(ctx context.Context, h host.Host) error {

	log.Debug("Starting DHT route discovery.")

	dhtInstance, err := initDHT(ctx, h)
	if err != nil {
		return err
	}

	routingDiscovery := drouting.NewRoutingDiscovery(dhtInstance)
	dutil.Advertise(ctx, routingDiscovery, ma.RENDEZVOUS)

	log.Infof("Starting DHT peer discovery for rendezvous string: %s", ma.RENDEZVOUS)

discoveryLoop:
	for {
		peerChan, err := routingDiscovery.FindPeers(ctx, ma.RENDEZVOUS)
		if err != nil {
			return fmt.Errorf("peer discovery error: %w", err)
		}

		for {
			select {
			case p, ok := <-peerChan:
				if !ok {
					peerChan = nil
					break
				}
				if p.ID == h.ID() {
					continue // Skip self connection
				}

				err := h.Connect(ctx, p) // Using the outer context directly
				if err != nil {
					log.Debugf("Failed connecting to %s, error: %v", p.ID.String(), err)
				} else {
					log.Infof("Connected to DHT peer: %s", p.ID.String())

					// Add peer to list of known peers
					peerMutex.Lock()
					connectedPeers[p.ID.String()] = &p
					peerMutex.Unlock()

					break discoveryLoop
				}
			case <-ctx.Done():
				log.Info("Context cancelled, stopping DHT peer discovery.")
				return nil
			}
			if peerChan == nil {
				break
			}
		}
	}

	log.Info("DHT Peer discovery complete")
	return nil
}
