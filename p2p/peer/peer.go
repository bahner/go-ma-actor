package peer

import (
	p2peer "github.com/libp2p/go-libp2p/core/peer"
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

func GetOrCreateFromAddrInfo(addrInfo *p2peer.AddrInfo) (Peer, error) {

	id := addrInfo.ID.String()

	p, err := Get(id)
	if err == nil {
		return p, nil
	}

	nodeAlias, err := LookupNick(id)
	if err != nil {
		nodeAlias = createNodeAlias(id)
	}

	return New(addrInfo, nodeAlias, defaultAllowed), nil

}

func createNodeAlias(id string) string {

	if len(id) <= defaultAliasLength {
		return id
	}

	return id[len(id)-defaultAliasLength:]

}
