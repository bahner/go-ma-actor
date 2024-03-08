package peer

import (
	"github.com/bahner/go-ma-actor/config"
	p2peer "github.com/libp2p/go-libp2p/core/peer"
)

const (
	DEFAULT_TAG_VALUE = 100
)

type Peer struct {
	// ID is the string representation of the peer's ID
	ID string
	// Name is the peer's name
	Nick string
	// AddrInfo
	AddrInfo *p2peer.AddrInfo
	// Allowed
	Allowed bool
}

// Create a new aliased addrinfo peer
func New(addrInfo *p2peer.AddrInfo, nick string, allowed bool) Peer {

	return Peer{
		ID:       addrInfo.ID.String(),
		AddrInfo: addrInfo,
		Nick:     nick,
		Allowed:  allowed,
	}
}

// Get or create a peer from an addrinfo. This is a dead function,
// in the sense that it does not do any live P2P lookups and as such
// it's use is safe to use anytime.
// The lookup is just in the local memory cache and database.
func GetOrCreateFromAddrInfo(addrInfo *p2peer.AddrInfo) (Peer, error) {

	id := addrInfo.ID.String()

	p, err := Get(id)
	if err == nil {
		return p, nil
	}

	nodeAlias, err := LookupNick(id)
	if err != nil {
		nodeAlias = addrInfo.ID.ShortString()
	}

	return New(addrInfo, nodeAlias, config.P2PDiscoveryAllowAll()), nil

}
