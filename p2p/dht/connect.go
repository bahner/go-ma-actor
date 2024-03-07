package dht

import (
	"context"

	"github.com/bahner/go-ma"
	"github.com/bahner/go-ma-actor/p2p/peer"
	log "github.com/sirupsen/logrus"
)

func (d *DHT) PeerConnectAndUpdateIfSuccessful(ctx context.Context, p peer.Peer) error {

	if len(p.AddrInfo.Addrs) == 0 {
		return ErrAddrInfoAddrsEmpty
	}

	id := p.AddrInfo.ID

	err := d.h.Connect(ctx, *p.AddrInfo)
	// NOOP. Clients that are protected are allowed to connect to us.
	// Even if we can't connect to them right now, we should still protect them.
	// if err != nil && d.h.ConnManager().IsProtected(id, ma.RENDEZVOUS) {
	// log.Warnf("Unprotecting previously protected peer %s: %v", id, err)
	// d.h.ConnManager().UntagPeer(id, ma.RENDEZVOUS)
	// d.h.ConnManager().Unprotect(id, ma.RENDEZVOUS)
	// }
	if err != nil {
		return err
	}

	if !d.h.ConnManager().IsProtected(id, ma.RENDEZVOUS) {
		log.Infof("Protecting previously unprotected peer %s", id)
		d.h.ConnManager().TagPeer(p.AddrInfo.ID, ma.RENDEZVOUS, peer.DEFAULT_TAG_VALUE)
		d.h.ConnManager().Protect(p.AddrInfo.ID, ma.RENDEZVOUS)

		// This is a new peer, so we should allow it explicitly.
		// ACtually it should be allowed by default, but we'll set it explicitly here.
		// Ref. line #99 above
		p.Allowed = true

	}

	return peer.Set(p)

}
