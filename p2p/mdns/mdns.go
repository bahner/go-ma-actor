package mdns

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	p2pmdns "github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	log "github.com/sirupsen/logrus"
)

type MDNS struct {
	PeerChan   chan peer.AddrInfo
	h          host.Host
	rendezvous string
}

// discoveryNotifee gets notified when a new peer is discovered.
type discoveryNotifee struct {
	PeerChan chan peer.AddrInfo
}

// HandlePeerFound is called when a new peer is discovered.
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	n.PeerChan <- pi
}

// New initializes the MDNS discovery service and returns an MDNS instance.
func New(h host.Host, rendezvous string) (*MDNS, error) {
	n := &discoveryNotifee{
		PeerChan: make(chan peer.AddrInfo),
	}

	// Initialize the MDNS service and start it
	service := p2pmdns.NewMdnsService(h, rendezvous, n)
	if err := service.Start(); err != nil {
		return nil, fmt.Errorf("failed to start MDNS service: %w", err)
	}

	// Since service is not stored, there's no direct reference to it in the MDNS struct
	return &MDNS{
		PeerChan:   n.PeerChan,
		h:          h,
		rendezvous: rendezvous,
	}, nil
}

// DiscoverPeers starts the discovery process and connects to discovered peers until the context is cancelled.
func (m *MDNS) DiscoverPeers(ctx context.Context) error {
	log.Debugf("Discovering MDNS peers for service name: %s", m.rendezvous)

	// Start trimming connections to make room for new peers
	m.h.ConnManager().TrimOpenConns(context.Background())

	for {
		select {
		case p, ok := <-m.PeerChan:
			if !ok {
				log.Debug("MDNS peer channel closed.")
				return nil // Exit if the channel is closed
			}
			if p.ID == m.h.ID() {
				continue // Skip self connection
			}

			if m.h.Network().Connectedness(p.ID) == network.Connected {
				log.Debugf("Already connected to MDNS peer: %s", p.ID.String())
				continue // Skip already connected peers
			}

			log.Infof("Found MDNS peer: %s, connecting", p.ID.String())
			err := m.h.Connect(ctx, p)
			if err != nil {
				log.Debugf("Failed connecting to %s, error: %v", p.ID.String(), err)
			} else {
				log.Infof("Connected to MDNS peer: %s", p.ID.String())
				// Add peer to list of known peers and protect the connection
				log.Debugf("Protecting discovered MDNS peer: %s", p.ID.String())
				m.h.ConnManager().TagPeer(p.ID, m.rendezvous, 10)
				m.h.ConnManager().Protect(p.ID, m.rendezvous)
				// Do not break; continue discovering
			}

		case <-ctx.Done():
			log.Info("Context cancelled, stopping MDNS peer discovery.")
			return nil // Stop the discovery loop if the context is done
		}
	}
}
