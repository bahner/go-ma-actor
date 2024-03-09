package peer

import (
	"errors"

	"github.com/bahner/go-ma-actor/config"
	"github.com/libp2p/go-libp2p/core/host"
	p2peer "github.com/libp2p/go-libp2p/core/peer"
)

const (
	DEFAULT_TAG_VALUE = 100
	defaultAllowAll   = true
	nodeAliasLength   = 8
)

var ErrIDTooShort = errors.New("ID too short")

type Peer struct {
	// ID is the string representation of the peer's ID
	ID string
	// Name is the peer's name
	Nick string
	// AddrInfo
	AddrInfo p2peer.AddrInfo
	// Allowed
	Allowed bool
}

// Create a new aliased addrinfo peer
func New(addrInfo p2peer.AddrInfo, nick string, allowed bool) Peer {

	return Peer{
		ID:       addrInfo.ID.String(),
		AddrInfo: addrInfo,
		Nick:     nick,
		Allowed:  allowed,
	}
}

// Get or creates a peer from the ID String.
// This might take sometime, but it's still very useful.
// It should normally e pretty fast.
func GetOrCreatePeerFromIDString(h host.Host, id string) (Peer, error) {

	_, err := p2peer.Decode(id)
	if err != nil {
		return Peer{}, err
	}

	_p, err := Get(id)
	if err == nil {
		// Always do a lookup on the nick as it might've changed
		_p.Nick = getOrCreateNick(id)
		return _p, nil
	}

	addrInfo, err := PeerAddrInfoFromPeerIDString(h, id)
	if err != nil {
		return Peer{}, err
	}

	return GetOrCreateFromAddrInfo(addrInfo)
}

// Get or create a peer from an addrinfo. This is a dead function,
// in the sense that it does not do any live P2P lookups and as such
// it's use is safe to use anytime.
// The lookup is just in the local memory cache and database.
func GetOrCreateFromAddrInfo(addrInfo p2peer.AddrInfo) (Peer, error) {

	id := addrInfo.ID.String()
	nick := getOrCreateNick(id)

	p, err := Get(id)
	if err == nil {
		p.Nick = nick
		return p, nil
	}

	return New(
		addrInfo,
		nick,
		config.DEFAULT_ALLOW_ALL), nil

}
