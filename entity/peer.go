package entity

import (
	"github.com/bahner/go-ma-actor/p2p"
	"github.com/bahner/go-ma-actor/p2p/peer"
	"github.com/libp2p/go-libp2p/core/network"
	p2peer "github.com/libp2p/go-libp2p/core/peer"
	log "github.com/sirupsen/logrus"
)

// This function connects to a peer using the DHT.
// The peer is identified by it's DID. The return value is the peer's AddrInfo.
// But you needn't use that for anything
func (e *Entity) ConnectPeer() (pi p2peer.AddrInfo, err error) {

	p := p2p.Get()

	pid, err := e.DID.PeerID()
	if err != nil {
		log.Debugf("Failed to get peer ID: %v", err)
		return p2peer.AddrInfo{}, err
	}

	// If we're already connected, return
	if p.Host.Network().Connectedness(pid) == network.Connected {
		log.Debugf("Already connected to peer: %s", pid.String())
		return pi, peer.ErrAlreadyConnected
	}

	// Look for the peer in the DHT
	pai, err := p.DHT.FindPeer(e.Ctx, pid)
	if err != nil {
		log.Debugf("Failed to find peer: %v", err)
		return pi, err
	}
	log.Debugf("PeerInfo: %v", pai.Addrs)

	// Connect to the peer
	log.Debugf("Connecting to peer with addrs: %v", pi.Addrs)
	err = p.Host.Connect(e.Ctx, pi)

	return pi, err

}
