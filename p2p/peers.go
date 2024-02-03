package p2p

import (
	"github.com/bahner/go-ma"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

// GetAllConnectedPeers returns a slice of peer.ID for all connected peers of the given host.
func (p *P2P) GetAllConnectedPeers() peer.IDSlice {
	h := p.Node
	var connectedPeers peer.IDSlice

	for _, p := range h.Network().Peers() {
		if h.Network().Connectedness(p) == network.Connected {
			connectedPeers = append(connectedPeers, p)
		}
	}

	return connectedPeers
}

// GetConnectedProtectedPeers returns a slice of peer.ID for all protected connected peers.
func (p *P2P) GetConnectedProtectedPeers() peer.IDSlice {
	h := p.Node
	var connectedProtectedPeers peer.IDSlice

	for _, connectedPeer := range p.GetAllConnectedPeers() {
		if h.ConnManager().IsProtected(connectedPeer, ma.RENDEZVOUS) {
			connectedProtectedPeers = append(connectedProtectedPeers, connectedPeer)
		}
	}

	return connectedProtectedPeers
}

// GetConnectedUnprotectedPeers returns a slice of peer.ID for all unprotected connected peers.
func (p *P2P) GetConnectedUnprotectedPeers() peer.IDSlice {
	connectedPeers := p.GetAllConnectedPeers()
	connectedProtectedPeers := p.GetConnectedProtectedPeers()

	var connectedUnprotectedPeers peer.IDSlice
	for _, connectedPeer := range connectedPeers {
		if !containsPeer(connectedProtectedPeers, connectedPeer) {
			connectedUnprotectedPeers = append(connectedUnprotectedPeers, connectedPeer)
		}
	}

	return connectedUnprotectedPeers
}

// containsPeer checks if a peer.ID is present in a slice of peer.ID.
func containsPeer(slice peer.IDSlice, peerID peer.ID) bool {
	for _, p := range slice {
		if p == peerID {
			return true
		}
	}
	return false
}

// GetConnectedProtectedPeersAddrInfo returns a map of peer.ID to AddrInfo for all protected connected peers.
func (p *P2P) GetConnectedProtectedPeersAddrInfo() map[string]*peer.AddrInfo {
	h := p.Node
	connectedPeersAddrInfo := make(map[string]*peer.AddrInfo)

	for _, connectedPeer := range p.GetConnectedProtectedPeers() {
		peerAddrInfo := h.Peerstore().PeerInfo(connectedPeer)
		connectedPeersAddrInfo[connectedPeer.String()] = &peerAddrInfo
	}

	return connectedPeersAddrInfo
}

func (p *P2P) GetConnectedProctectedPeersShortStrings() []string {
	peers := p.GetConnectedProtectedPeersAddrInfo()
	peersShortstrings := make([]string, 0, len(peers))
	for _, p := range peers {
		peersShortstrings = append(peersShortstrings, p.ID.ShortString())
	}
	return peersShortstrings
}
