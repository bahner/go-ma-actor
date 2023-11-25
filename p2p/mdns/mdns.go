package mdns

import (
	"context"

	"github.com/bahner/go-ma"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"

	p2pmdns "github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	log "github.com/sirupsen/logrus"
)

type discoveryNotifee struct {
	PeerChan chan peer.AddrInfo
}

// interface to be called when new  peer is found
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	n.PeerChan <- pi
}

// Initialize the MDNS service
func initMDNS(h host.Host, rendezvous string) chan peer.AddrInfo {
	// register with service so that we get notified about peer discovery
	n := &discoveryNotifee{}
	n.PeerChan = make(chan peer.AddrInfo)

	// An hour might be a long long period in practical applications. But this is fine for us
	ser := p2pmdns.NewMdnsService(h, rendezvous, n)
	if err := ser.Start(); err != nil {
		panic(err)
	}
	return n.PeerChan
}
func DiscoverPeers(ctx context.Context, h host.Host) error {
	log.Debugf("Discovering MDNS peers for servicename: %s", ma.RENDEZVOUS)

	peerChan := initMDNS(h, ma.RENDEZVOUS)
	// Start trimming connections, so we have room for new friends
	h.ConnManager().TrimOpenConns(context.Background())

discoveryLoop:
	for {
		select {
		case p, ok := <-peerChan:
			if !ok {
				log.Debug("MDNS peer channel closed.")
				break discoveryLoop
			}
			if p.ID == h.ID() {
				continue // Skip self connection
			}

			log.Infof("Found MDNS peer: %s connecting", p.ID.String())
			err := h.Connect(ctx, p)
			if err != nil {
				log.Debugf("Failed connecting to %s, error: %v", p.ID.String(), err)
			} else {
				log.Infof("Connected to MDNS peer: %s", p.ID.String())

				// Add peer to list of known peers
				log.Debugf("Protecting discovered MDNS peer: %s", p.ID.String())
				h.ConnManager().TagPeer(p.ID, ma.RENDEZVOUS, 10)
				h.ConnManager().Protect(p.ID, ma.RENDEZVOUS)

			}

		case <-ctx.Done():
			log.Info("Context cancelled, stopping MDNS peer discovery.")
			return nil
		}
	}

	log.Info("MDNS peer discovery complete")
	return nil
}
