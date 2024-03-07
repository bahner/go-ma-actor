package mdns

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	p2pmdns "github.com/libp2p/go-libp2p/p2p/discovery/mdns"
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
