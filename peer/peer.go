package peer

import (
	"github.com/bahner/go-ma-actor/alias"
	p2peer "github.com/libp2p/go-libp2p/core/peer"
)

type Peer struct {
	// ID is the peer's ID
	ID string
	// Name is the peer's name
	Alias string
	// AddrInfo
	AddrInfo *p2peer.AddrInfo
}

// Create a new aliased addrinfo peer
func NewWithAlias(addrInfo *p2peer.AddrInfo, alias string) *Peer {

	id := addrInfo.ID.String()
	return &Peer{
		ID:       id,
		Alias:    alias,
		AddrInfo: addrInfo,
	}
}

// Create a new aliased addrinfo peer
func New(addrInfo *p2peer.AddrInfo) *Peer {
	na := alias.GetNodeAlias(addrInfo.ID.String())
	if na == "" {
		na = addrInfo.ID.String()
		na = na[len(na)-8:]
	}
	return NewWithAlias(addrInfo, na)
}

func GetOrCreate(addrInfo *p2peer.AddrInfo) *Peer {

	id := addrInfo.ID.String()

	p := get(id)
	if p == nil {
		p = New(addrInfo)
		add(p)
	}
	return p
}
