package peer

import (
	"context"

	"github.com/bahner/go-ma"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	p2peer "github.com/libp2p/go-libp2p/core/peer"
	log "github.com/sirupsen/logrus"
)

// Get or creates a peer from the ID String.
// This might take sometime, but it's still very useful.
// It should normally e pretty fast.
func GetOrCreatePeerFromIDString(h host.Host, id string) (Peer, error) {

	_p, err := Get(id)
	if err == nil {
		return _p, nil
	}

	addrInfo, err := GetPeerAddrInfoFromIDString(h, id)
	if err != nil {
		return Peer{}, err
	}

	return GetOrCreateFromAddrInfo(addrInfo)
}

func GetPeerAddrInfoFromIDString(h host.Host, id string) (p2peer.AddrInfo, error) {
	pid, err := p2peer.Decode(id)
	if err != nil {
		return p2peer.AddrInfo{}, err
	}

	return GetPeerAddrInfoFromID(h, pid)
}

func GetPeerAddrInfoFromID(h host.Host, id p2peer.ID) (p2peer.AddrInfo, error) {
	a := p2peer.AddrInfo{
		ID:    id,
		Addrs: h.Peerstore().Addrs(id),
	}

	return a, nil
}

func ConnectAndProtect(ctx context.Context, h host.Host, pai p2peer.AddrInfo) error {

	var (
		p  Peer
		id = pai.ID.String()
	)
	if len(pai.Addrs) == 0 {
		return ErrAddrInfoAddrsEmpty
	}

	p, err := GetOrCreateFromAddrInfo(pai)
	if err != nil {
		return err
	}
	if !IsAllowed(p.ID) { // Do an actual lookup in the database here
		log.Debugf("Peer %s is explicitly denied", id)
		UnprotectPeer(h, pai.ID)
		return ErrPeerDenied
	}

	if h.Network().Connectedness(pai.ID) == network.Connected {
		log.Debugf("Already connected to DHT peer: %s", id)
		return ErrAlreadyConnected // This is not an error, but we'll return it as such for now.
	}

	err = Protect(h, pai.ID)
	if err != nil {
		log.Warnf("Failed to protect peer %s: %v", id, err)
	}

	err = h.Connect(ctx, pai)
	if err != nil {
		log.Warnf("Failed to connect to peer %s: %v", id, err)
		return err
	}

	return Set(p)
}

func Protect(h host.Host, id p2peer.ID) error {

	if !h.ConnManager().IsProtected(id, ma.RENDEZVOUS) {
		log.Infof("Protecting previously unprotected peer %s", id.String())
		h.ConnManager().TagPeer(id, ma.RENDEZVOUS, DEFAULT_TAG_VALUE)
		h.ConnManager().Protect(id, ma.RENDEZVOUS)
	}

	return nil
}

func UnprotectPeer(h host.Host, id p2peer.ID) error {

	if h.ConnManager().IsProtected(id, ma.RENDEZVOUS) {
		log.Infof("Unprotecting previously protected peer %s", id.String())
		h.ConnManager().UntagPeer(id, ma.RENDEZVOUS)
		h.ConnManager().Unprotect(id, ma.RENDEZVOUS)
	}

	return nil
}
