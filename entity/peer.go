package entity

import (
	"github.com/bahner/go-ma-actor/p2p"
	"github.com/libp2p/go-libp2p/core/network"
)

// This function connects to a peer using the DHT.
// The peer is identified by it's DID.
func (e *Entity) ConnectPeer() error {

	p := p2p.Get()

	pid, err := e.DID.PeerID()
	if err != nil {
		return err
	}

	// If we're already connected, return
	if p.DHT.Host().Network().Connectedness(pid) == network.Connected {
		return nil
	}

	// Look for the peer in the DHT
	pi, err := p.DHT.FindPeer(e.Ctx, pid)
	if err != nil {
		return err
	}

	// Connect to the peer
	return p.DHT.Host().Connect(e.Ctx, pi)

}
