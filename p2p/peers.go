package p2p

import (
	"github.com/bahner/go-ma"
	"github.com/bahner/go-ma-actor/p2p/peer"
	"github.com/libp2p/go-libp2p/core/network"
	p2peer "github.com/libp2p/go-libp2p/core/peer"
)

// AllConnectedPeers returns a slice of p2peer.ID for all connected peers of the given host.
func (p *P2P) AllConnectedPeers() p2peer.IDSlice {
	h := p.Host
	var connectedPeers p2peer.IDSlice

	for _, p := range h.Network().Peers() {
		if h.Network().Connectedness(p) == network.Connected {
			connectedPeers = append(connectedPeers, p)
		}
	}

	return connectedPeers
}

// ConnectedProtectedPeers returns a slice of p2peer.ID for all protected connected peers.
func (p *P2P) ConnectedProtectedPeers() p2peer.IDSlice {
	h := p.Host
	var connectedProtectedPeers p2peer.IDSlice

	for _, connectedPeer := range p.AllConnectedPeers() {
		if h.ConnManager().IsProtected(connectedPeer, ma.RENDEZVOUS) {
			connectedProtectedPeers = append(connectedProtectedPeers, connectedPeer)
		}
	}

	return connectedProtectedPeers
}

// ConnectedUnprotectedPeers returns a slice of p2peer.ID for all unprotected connected peers.
func (p *P2P) ConnectedUnprotectedPeers() p2peer.IDSlice {
	h := p.Host
	var connectedUnprotectedPeers p2peer.IDSlice

	for _, connectedPeer := range p.AllConnectedPeers() {
		if !h.ConnManager().IsProtected(connectedPeer, ma.RENDEZVOUS) {
			connectedUnprotectedPeers = append(connectedUnprotectedPeers, connectedPeer)
		}
	}

	return connectedUnprotectedPeers
}

// ConnectedProtectedPeersAddrInfo returns a map of p2peer.ID to AddrInfo for all protected connected peers.
func (p *P2P) ConnectedProtectedPeersAddrInfo() map[string]p2peer.AddrInfo {
	h := p.Host
	connectedPeersAddrInfo := make(map[string]p2peer.AddrInfo)

	for _, connectedPeer := range p.ConnectedProtectedPeers() {
		peerAddrInfo := h.Peerstore().PeerInfo(connectedPeer)
		connectedPeersAddrInfo[connectedPeer.String()] = peerAddrInfo
	}

	return connectedPeersAddrInfo
}

func (p *P2P) ConnectedProctectedPeersNickList() []string {
	peers := p.ConnectedProtectedPeersAddrInfo()
	peersNickList := make([]string, 0, len(peers))
	for _, p := range peers {
		nick, err := peer.Nick(p.ID.String())
		if err != nil {
			// We only want known live peers in the list
			continue
		}
		peersNickList = append(peersNickList, nick)
	}
	return peersNickList
}
