package mdns

import (
	"context"

	"github.com/bahner/go-ma"
	"github.com/bahner/go-ma-actor/p2p/peer"
	"github.com/libp2p/go-libp2p/core/network"
	p2peer "github.com/libp2p/go-libp2p/core/peer"
	log "github.com/sirupsen/logrus"
)

func (m *MDNS) peerConnectAndUpdateIfSuccessful(ctx context.Context, pai p2peer.AddrInfo) error {

	var p peer.Peer

	if len(pai.Addrs) == 0 {
		return peer.ErrAddrInfoAddrsEmpty
	}
	if m.h.Network().Connectedness(pai.ID) == network.Connected {
		log.Debugf("Already connected to MDNS peer: %s", pai.ID.String())
		return peer.ErrAlreadyConnected // This is not an error, but we'll return it as such for now.
	}

	err := m.h.Connect(ctx, pai)
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

	if !m.h.ConnManager().IsProtected(pai.ID, ma.RENDEZVOUS) {
		log.Infof("Protecting previously unprotected peer %s", pai.ID.String())
		m.h.ConnManager().TagPeer(pai.ID, ma.RENDEZVOUS, peer.DEFAULT_TAG_VALUE)
		m.h.ConnManager().Protect(pai.ID, ma.RENDEZVOUS)

		// This is a new peer, so we should allow it explicitly.
		// ACtually it should be allowed by default, but we'll set it explicitly here.
		// Ref. line #99 above

		p, err = peer.GetOrCreateFromAddrInfo(&pai)
		if err != nil {
			return err
		}
		p.Allowed = true

	}

	return peer.Set(p)

}
