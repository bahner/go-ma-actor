package entity

import (
	"github.com/bahner/go-ma-actor/p2p"
	"github.com/bahner/go-ma-actor/p2p/peer"
	"github.com/libp2p/go-libp2p/core/network"
	p2peer "github.com/libp2p/go-libp2p/core/peer"
	log "github.com/sirupsen/logrus"
)

// This function connects to a peer using the DHT.
// The peer is identified by it's DID. The return value is the peer's AddrInfo,
// which is likely only useful for debugging and information.
func (e *Entity) ConnectPeer() (pai p2peer.AddrInfo, err error) {

	p := p2p.Get()
	pid := e.DID.Name.Peer()

	// If we're already connected, return
	if p.Host.Network().Connectedness(pid) == network.Connected {
		log.Debugf("Already connected to peer: %s", pid.String())
		return p.Host.Peerstore().PeerInfo(pid), peer.ErrAlreadyConnected
	}

	// Look for the peer in the DHT
	pai, err = p.DHT.FindPeer(e.Ctx, pid)
	if err != nil {
		log.Debugf("Failed to find peer: %v", err)
		// return pi, err
	}
	log.Debugf("PeerInfo: %v", pai.Addrs)

	// Connect to the peer
	log.Debugf("Connecting to peer with addrs: %v", pai.Addrs)
	err = p.Host.Connect(e.Ctx, pai)

	return pai, err

}
