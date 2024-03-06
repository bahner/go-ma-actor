package p2p

import (
	"context"

	"github.com/bahner/go-ma"
	"github.com/bahner/go-ma-actor/p2p/peer"
	"github.com/libp2p/go-libp2p/core/network"
	p2peer "github.com/libp2p/go-libp2p/core/peer"
)

// GetAllConnectedPeers returns a slice of p2peer.ID for all connected peers of the given host.
func (p *P2P) GetAllConnectedPeers() p2peer.IDSlice {
	h := p.Node
	var connectedPeers p2peer.IDSlice

	for _, p := range h.Network().Peers() {
		if h.Network().Connectedness(p) == network.Connected {
			connectedPeers = append(connectedPeers, p)
		}
	}

	return connectedPeers
}

// GetConnectedProtectedPeers returns a slice of p2peer.ID for all protected connected peers.
func (p *P2P) GetConnectedProtectedPeers() p2peer.IDSlice {
	h := p.Node
	var connectedProtectedPeers p2peer.IDSlice

	for _, connectedPeer := range p.GetAllConnectedPeers() {
		if h.ConnManager().IsProtected(connectedPeer, ma.RENDEZVOUS) {
			connectedProtectedPeers = append(connectedProtectedPeers, connectedPeer)
		}
	}

	return connectedProtectedPeers
}

// GetConnectedUnprotectedPeers returns a slice of p2peer.ID for all unprotected connected peers.
func (p *P2P) GetConnectedUnprotectedPeers() p2peer.IDSlice {
	connectedPeers := p.GetAllConnectedPeers()
	connectedProtectedPeers := p.GetConnectedProtectedPeers()

	var connectedUnprotectedPeers p2peer.IDSlice
	for _, connectedPeer := range connectedPeers {
		if !containsPeer(connectedProtectedPeers, connectedPeer) {
			connectedUnprotectedPeers = append(connectedUnprotectedPeers, connectedPeer)
		}
	}

	return connectedUnprotectedPeers
}

// containsPeer checks if a p2peer.ID is present in a slice of p2peer.ID.
func containsPeer(slice p2peer.IDSlice, peerID p2peer.ID) bool {
	for _, p := range slice {
		if p == peerID {
			return true
		}
	}
	return false
}

// GetConnectedProtectedPeersAddrInfo returns a map of p2peer.ID to AddrInfo for all protected connected peers.
func (p *P2P) GetConnectedProtectedPeersAddrInfo() map[string]*p2peer.AddrInfo {
	h := p.Node
	connectedPeersAddrInfo := make(map[string]*p2peer.AddrInfo)

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

// Get or creates a peer from the ID
// NB! This is a heavy operation and should be used with caution
func (_p *P2P) GetOrCreatePeerFromIDString(id string) (peer.Peer, error) {

	p, err := peer.Get(id)
	if err == nil {
		return p, nil
	}

	pid, err := p2peer.Decode(id)
	if err != nil {
		return peer.Peer{}, err
	}

	addrInfo, err := _p.DHT.FindPeer(context.Background(), pid)
	if err != nil {
		return peer.Peer{}, err
	}

	p = peer.New(&addrInfo)
	err = peer.Set(p)
	if err != nil {
		return peer.Peer{}, err
	}

	return p, nil
}
