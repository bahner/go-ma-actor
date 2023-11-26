package p2p

import (
	"github.com/bahner/go-ma"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

var (
	connectedPeers = make(map[string]*peer.AddrInfo)
)

// Get list of connected peers for the given host
func GetConnectedPeers(h host.Host) map[string]*peer.AddrInfo {

	for _, p := range h.Network().Peers() {

		if h.ConnManager().IsProtected(p, ma.RENDEZVOUS) {

			if h.Network().Connectedness(p) == network.Connected {

				connectedPeer := h.Peerstore().PeerInfo(p)

				connectedPeers[p.String()] = &connectedPeer
			}
		}

	}

	return connectedPeers
}
