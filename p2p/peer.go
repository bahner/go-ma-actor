package p2p

import (
	"github.com/bahner/go-ma-actor/p2p/peer"
	p2peer "github.com/libp2p/go-libp2p/core/peer"
)

// Get or creates a peer from the ID String.
// This might take sometime, but it's still very useful.
// It should normally e pretty fast.
func (p *P2P) GetOrCreatePeerFromIDString(id string) (peer.Peer, error) {

	_p, err := peer.Get(id)
	if err == nil {
		return _p, nil
	}

	addrInfo, err := p.GetPeerAddrInfoFromIDString(id)
	if err != nil {
		return peer.Peer{}, err
	}

	return peer.GetOrCreateFromAddrInfo(addrInfo)
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
