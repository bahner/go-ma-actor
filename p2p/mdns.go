package p2p

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma"
	"github.com/bahner/go-ma-actor/p2p/peer"
	log "github.com/sirupsen/logrus"

	"github.com/libp2p/go-libp2p/core/host"
	p2peer "github.com/libp2p/go-libp2p/core/peer"
	p2pmdns "github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

type MDNS struct {
	PeerChan chan p2peer.AddrInfo
	h        host.Host
}

// discoveryNotifee gets notified when a new peer is discovered.
type discoveryNotifee struct {
	PeerChan chan p2peer.AddrInfo
}

// HandlePeerFound is called when a new peer is discovered.
func (n *discoveryNotifee) HandlePeerFound(pi p2peer.AddrInfo) {
	n.PeerChan <- pi
}

// newMDNS initializes the MDNS discovery service and returns an MDNS instance.
func newMDNS(h host.Host) (*MDNS, error) {
	n := &discoveryNotifee{
		PeerChan: make(chan p2peer.AddrInfo),
	}

	// Initialize the MDNS service and start it
	service := p2pmdns.NewMdnsService(h, ma.RENDEZVOUS, n)
	if err := service.Start(); err != nil {
		return nil, fmt.Errorf("failed to start MDNS service: %w", err)
	}

	// Since service is not stored, there's no direct reference to it in the MDNS struct
	return &MDNS{
		PeerChan: n.PeerChan,
		h:        h,
	}, nil
}

// DiscoveryLoop starts the discovery process and connects to discovered peers until the context is cancelled.
func (m *MDNS) discoveryLoop(ctx context.Context) error {
	log.Debugf("Discovering MDNS peers for service name: %s", ma.RENDEZVOUS)

	for {
		select {
		case pai, ok := <-m.PeerChan:
			if !ok {
				if !(ctx == nil) { // conext.Bacground() is nil
					log.Fatalf("MDNS peer channel closed, ut was supposed to be running in the background.")
				}
				log.Debug("MDNS peer channel closed.")
				return nil // Exit if the channel is closed
			}

			if pai.ID == m.h.ID() {
				continue // Skip self connection
			}

			err := peer.ConnectAndProtect(ctx, m.h, pai)
			if err != nil {
				log.Debugf("Failed connecting to MDNS peer: %s, error: %v", pai.ID, err)
				continue // Skip if connection failed
			}

			log.Info("Successfully connected to MDNS peer: ", pai.ID)

		case <-ctx.Done():
			log.Info("Context cancelled, stopping MDNS peer discovery.")
			return nil // Stop the discovery loop if the context is done
		}
	}
}
