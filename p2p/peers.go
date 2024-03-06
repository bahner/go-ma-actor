package p2p

import (
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
	h := p.Node
	var connectedUnprotectedPeers p2peer.IDSlice

	for _, connectedPeer := range p.GetAllConnectedPeers() {
		if !h.ConnManager().IsProtected(connectedPeer, ma.RENDEZVOUS) {
			connectedUnprotectedPeers = append(connectedUnprotectedPeers, connectedPeer)
		}
	}

	return connectedUnprotectedPeers
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

func (p *P2P) GetPeerAddrInfoFromIDString(id string) (*p2peer.AddrInfo, error) {
	pid, err := p2peer.Decode(id)
	if err != nil {
		return nil, err
	}

	return p.GetPeerAddrInfoFromID(pid)
}

func (p *P2P) GetPeerAddrInfoFromID(id p2peer.ID) (*p2peer.AddrInfo, error) {
	a := p2peer.AddrInfo{
		ID:    id,
		Addrs: p.Node.Peerstore().Addrs(id),
	}

	return &a, nil
}

// Get or creates a peer from the ID
// NB! This is a heavy operation and should be used with caution
func (p *P2P) GetOrCreatePeerFromIDString(id string) (peer.Peer, error) {

	addrInfo, err := p.GetPeerAddrInfoFromIDString(id)
	if err != nil {
		return peer.Peer{}, err
	}

	return peer.GetOrCreateFromAddrInfo(addrInfo)
}
