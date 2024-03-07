package entity

import (
	"github.com/bahner/go-ma-actor/p2p"
	"github.com/bahner/go-ma-actor/p2p/peer"
	"github.com/libp2p/go-libp2p/core/network"
	p2peer "github.com/libp2p/go-libp2p/core/peer"
)

// This function connects to a peer using the DHT.
// The peer is identified by it's DID. The return value is the peer's AddrInfo.
// But you needn't use that for anything
func (e *Entity) ConnectPeer() (pi p2peer.AddrInfo, err error) {

	p := p2p.Get()

	pid, err := e.DID.PeerID()
	if err != nil {
		return p2peer.AddrInfo{}, err
	}

	// If we're already connected, return
	if p.DHT.Host().Network().Connectedness(pid) == network.Connected {
		return p2peer.AddrInfo{}, peer.ErrAlreadyConnected
	}

	// Look for the peer in the DHT
	pi, err = p.DHT.FindPeer(e.Ctx, pid)
	if err != nil {
		return pi, err
	}

	// Connect to the peer
	err = p.DHT.Host().Connect(e.Ctx, pi)

	return pi, err

}
