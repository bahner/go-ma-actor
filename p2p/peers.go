package p2p

import (
	"time"

	"github.com/bahner/go-ma"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

var (
	connectedPeers = make(map[string]*peer.AddrInfo)
)

// Get list of connected peers.
func GetConnectedPeers(connectTimeout time.Duration) map[string]*peer.AddrInfo {

	for _, p := range n.Network().Peers() {

		if n.ConnManager().IsProtected(p, ma.RENDEZVOUS) {

			if n.Network().Connectedness(p) == network.Connected {

				connectedPeer := n.Peerstore().PeerInfo(p)

				connectedPeers[p.String()] = &connectedPeer
			}
		}

	}

	return connectedPeers
}
