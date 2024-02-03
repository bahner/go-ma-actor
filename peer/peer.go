package peer

import (
	"fmt"

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
func New(addrInfo *p2peer.AddrInfo, alias string) *Peer {

	id := addrInfo.ID.String()
	return &Peer{
		ID:       id,
		Alias:    alias,
		AddrInfo: addrInfo,
	}
}

func GetOrCreate(addrInfo *p2peer.AddrInfo) (*Peer, error) {

	id := addrInfo.ID.String()

	na, err := alias.GetOrCreateNodeAlias(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get or create node alias: %w", err)
	}

	p := get(id)
	if p == nil {
		p = New(addrInfo, na)
		add(p)
	}

	return p, nil
}
