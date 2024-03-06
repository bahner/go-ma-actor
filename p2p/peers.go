package p2p

import (
	"github.com/bahner/go-ma"
	"github.com/libp2p/go-libp2p/core/network"
	p2peer "github.com/libp2p/go-libp2p/core/peer"
)

// AllConnectedPeers returns a slice of p2peer.ID for all connected peers of the given host.
func (p *P2P) AllConnectedPeers() p2peer.IDSlice {
	h := p.Node
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
	h := p.Node
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
	h := p.Node
	var connectedUnprotectedPeers p2peer.IDSlice

	for _, connectedPeer := range p.AllConnectedPeers() {
		if !h.ConnManager().IsProtected(connectedPeer, ma.RENDEZVOUS) {
			connectedUnprotectedPeers = append(connectedUnprotectedPeers, connectedPeer)
		}
	}

	return connectedUnprotectedPeers
}

// ConnectedProtectedPeersAddrInfo returns a map of p2peer.ID to AddrInfo for all protected connected peers.
func (p *P2P) ConnectedProtectedPeersAddrInfo() map[string]*p2peer.AddrInfo {
	h := p.Node
	connectedPeersAddrInfo := make(map[string]*p2peer.AddrInfo)

	for _, connectedPeer := range p.ConnectedProtectedPeers() {
		peerAddrInfo := h.Peerstore().PeerInfo(connectedPeer)
		connectedPeersAddrInfo[connectedPeer.String()] = &peerAddrInfo
	}

	return connectedPeersAddrInfo
}

func (p *P2P) ConnectedProctectedPeersShortStrings() []string {
	peers := p.ConnectedProtectedPeersAddrInfo()
	peersShortstrings := make([]string, 0, len(peers))
	for _, p := range peers {
		peersShortstrings = append(peersShortstrings, p.ID.ShortString())
	}
	return peersShortstrings
}
