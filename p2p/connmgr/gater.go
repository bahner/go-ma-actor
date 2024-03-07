package connmgr

import (
	"github.com/bahner/go-ma"
	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/p2p/peer"
	"github.com/libp2p/go-libp2p/core/control"
	"github.com/libp2p/go-libp2p/core/network"
	p2peer "github.com/libp2p/go-libp2p/core/peer"
	p2pConnmgr "github.com/libp2p/go-libp2p/p2p/net/connmgr"
	"github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"
)

// ConnectionGater is a struct that implements the network.ConnectionGater interface.
// It uses a sync.Map to store valid peer IDs that have been discovered using the correct rendezvous string.
type ConnectionGater struct {
	AllowAll bool
	ConnMgr  *p2pConnmgr.BasicConnMgr
}

// New creates a new CustomConnectionGater instance.
func NewConnectionGater(connMgr *p2pConnmgr.BasicConnMgr) *ConnectionGater {
	return &ConnectionGater{
		AllowAll: config.P2PDiscoveryAllow(),
		ConnMgr:  connMgr,
	}
}

// InterceptPeerDial checks if we should allow dialing the specified peer.
func (cg *ConnectionGater) InterceptPeerDial(p p2peer.ID) (allow bool) {

	return true
}

// InterceptAccept checks if an incoming connection from the specified network address should be allowed.
func (cg *ConnectionGater) InterceptAccept(conn network.ConnMultiaddrs) (allow bool) {
	return true
}

// InterceptSecured, InterceptUpgraded, and other methods can be implemented as needed.
// For simplicity, they are set to allow all connections in this example.
func (cg *ConnectionGater) InterceptSecured(nd network.Direction, p p2peer.ID, _ network.ConnMultiaddrs) (allow bool) {

	if nd == network.DirOutbound {
		return true
	}

	allow = cg.isAllowed(p)

	if allow {
		log.Debugf("InterceptSecured: Allow dialing to %s", p)
	} else {
		log.Debugf("InterceptSecured: Block dialing to %s", p)
	}
	return allow
}

func (cg *ConnectionGater) InterceptAddrDial(p p2peer.ID, _ multiaddr.Multiaddr) (allow bool) {

	return true
}

func (cg *ConnectionGater) InterceptUpgraded(_ network.Conn) (allow bool, reason control.DisconnectReason) {
	return true, 0
}

func (cg *ConnectionGater) isAllowed(p p2peer.ID) bool {

	if config.P2PDiscoveryAllow() {
		return true
	}

	// NB! Check peer.IsAllowed first. Because it might be explicitly denied and we want to adhere to that.
	// So if it's explicitly denied, we don't need to check the other conditions.
	if !peer.IsAllowed(p.String()) {
		log.Warnf("Peer %s is explicitly denied", p)
		return false
	}

	return cg.AllowAll || cg.ConnMgr.IsProtected(p, ma.RENDEZVOUS)
}
