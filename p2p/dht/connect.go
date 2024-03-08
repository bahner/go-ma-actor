package dht

import (
	"context"

	"github.com/bahner/go-ma"
	"github.com/bahner/go-ma-actor/p2p/peer"
	"github.com/libp2p/go-libp2p/core/network"
	p2peer "github.com/libp2p/go-libp2p/core/peer"
	log "github.com/sirupsen/logrus"
)

func (d *DHT) PeerConnectAndUpdateIfSuccessful(ctx context.Context, pai p2peer.AddrInfo) error {

	var p peer.Peer

	if len(pai.Addrs) == 0 {
		return peer.ErrAddrInfoAddrsEmpty
	}

	p, err := peer.GetOrCreateFromAddrInfo(&pai)
	if err != nil {
		return err
	}
	if !peer.IsAllowed(p.ID) { // Do an actual lookup in the database here
		log.Debugf("Peer %s is explicitly denied", pai.ID.String())
		d.unprotectPeer(pai.ID)
		return peer.ErrPeerDenied
	}

	if d.h.Network().Connectedness(pai.ID) == network.Connected {
		log.Debugf("Already connected to DHT peer: %s", pai.ID.String())
		return peer.ErrAlreadyConnected // This is not an error, but we'll return it as such for now.
	}

	err = d.protectPeer(pai.ID)
	if err != nil {
		log.Warnf("Failed to protect peer %s: %v", pai.ID.String(), err)
	}

	err = d.h.Connect(ctx, pai)
	if err != nil {
		return err
	}

	return peer.Set(p)
}

func (d *DHT) protectPeer(id p2peer.ID) error {

	if !d.h.ConnManager().IsProtected(id, ma.RENDEZVOUS) {
		log.Infof("Protecting previously unprotected peer %s", id.String())
		d.h.ConnManager().TagPeer(id, ma.RENDEZVOUS, peer.DEFAULT_TAG_VALUE)
		d.h.ConnManager().Protect(id, ma.RENDEZVOUS)
	}

	return nil
}

func (d *DHT) unprotectPeer(id p2peer.ID) error {

	if d.h.ConnManager().IsProtected(id, ma.RENDEZVOUS) {
		log.Infof("Unprotecting previously protected peer %s", id.String())
		d.h.ConnManager().UntagPeer(id, ma.RENDEZVOUS)
		d.h.ConnManager().Unprotect(id, ma.RENDEZVOUS)
	}

	return nil
}
