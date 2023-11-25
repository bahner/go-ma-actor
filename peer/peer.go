package peer

import p2peer "github.com/libp2p/go-libp2p/core/peer"

type Peer struct {
	// ID is the peer's ID
	ID string
	// Name is the peer's name
	Alias string
	// AddrInfo
	AddrInfo *p2peer.AddrInfo
}

func NewWithAlias(addrInfo *p2peer.AddrInfo, alias string) *Peer {

	id := addrInfo.ID.String()
	return &Peer{
		ID:       id,
		Alias:    alias,
		AddrInfo: addrInfo,
	}
}

func New(addrInfo *p2peer.AddrInfo) *Peer {
	alias := addrInfo.ID.String()
	return NewWithAlias(addrInfo, alias[len(alias)-8:])
}

func GetOrCreate(addrInfo *p2peer.AddrInfo) *Peer {

	id := addrInfo.ID.String()

	p := Get(id)
	if p == nil {
		p = New(addrInfo)
		Add(p)
	}
	return p
}
