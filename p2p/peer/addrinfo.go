package peer

import (
	"github.com/libp2p/go-libp2p/core/host"
	p2peer "github.com/libp2p/go-libp2p/core/peer"
)

func PeerAddrInfoFromPeerIDString(h host.Host, id string) (p2peer.AddrInfo, error) {
	pid, err := p2peer.Decode(id)
	if err != nil {
		return p2peer.AddrInfo{}, err
	}

	return PeerAddrInfoFromID(h, pid)
}

func PeerAddrInfoFromID(h host.Host, id p2peer.ID) (p2peer.AddrInfo, error) {
	a := p2peer.AddrInfo{
		ID:    id,
		Addrs: h.Peerstore().Addrs(id),
	}

	return a, nil
}
