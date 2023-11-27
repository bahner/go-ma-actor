package dht

import (
	"context"
	"fmt"
	"sync"

	"github.com/bahner/go-ma"
	p2pDHT "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
	log "github.com/sirupsen/logrus"
)

var dhtStructInstance *dhtStruct
var _ DHT = (*dhtStruct)(nil)

type DHT interface {
	DiscoverPeers(context.Context) error
}

type dhtStruct struct {
	*p2pDHT.IpfsDHT
	h host.Host
}

func Init(ctx context.Context, h host.Host, dhtOpts ...p2pDHT.Option) (DHT, error) {
	log.Info("Initializing Kademlia DHT.")

	var err error
	dhtStructInstance = &dhtStruct{h: h}

	dhtStructInstance.IpfsDHT, err = p2pDHT.New(ctx, h, dhtOpts...)
	if err != nil {
		log.Error("Failed to create Kademlia DHT.")
		return nil, err
	} else {
		log.Debug("Kademlia DHT created.")
	}

	err = dhtStructInstance.IpfsDHT.Bootstrap(ctx)
	if err != nil {
		log.Error("Failed to bootstrap Kademlia DHT.")
		return nil, err
	} else {
		log.Debug("Kademlia DHT bootstrap setup.")
	}

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
	return dhtStructInstance, nil
}

// Takes a context and a DHT instance and discovers peers using the DHT.
// You might want to se server option or not for the DHT.
func (d *dhtStruct) DiscoverPeers(ctx context.Context) error {

	log.Debug("Starting DHT route discovery.")

	routingDiscovery := drouting.NewRoutingDiscovery(d.IpfsDHT)
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
				if p.ID == d.h.ID() {
					continue // Skip self connection
				}

				err := d.h.Connect(ctx, p) // Using the outer context directly
				if err != nil {
					log.Debugf("Failed connecting to %s, error: %v", p.ID.String(), err)
				} else {
					log.Infof("Connected to DHT peer: %s", p.ID.String())

					// Add peer to list of known peers
					log.Debugf("Protecting peer: %s", p.ID.String())
					d.h.ConnManager().TagPeer(p.ID, ma.RENDEZVOUS, 10)
					d.h.ConnManager().Protect(p.ID, ma.RENDEZVOUS)

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

func Get() *p2pDHT.IpfsDHT {
	return dhtStructInstance.IpfsDHT
}
