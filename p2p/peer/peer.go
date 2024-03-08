package peer

import (
	"database/sql"
	"errors"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/config/db"
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

// Get or create a peer from an addrinfo. This is a dead function,
// in the sense that it does not do any live P2P lookups and as such
// it's use is safe to use anytime.
// The lookup is just in the local memory cache and database.
func GetOrCreateFromAddrInfo(addrInfo p2peer.AddrInfo) (Peer, error) {

	id := addrInfo.ID.String()

	p, err := Get(id)
	if err == nil {
		return p, nil
	}

	return New(
		addrInfo,
		getOrCreateNodeAlias(id),
		config.DEFAULT_ALLOW_ALL), nil

}

// Return a boolean whther the peer is knoor not
// This this should err on the side of caution and return false
func IsKnown(id string) bool {
	db, err := db.Get()
	if err != nil {
		return false
	}

	// We just need to know if the peer exists, so we select the id itself.
	var peerID string
	err = db.QueryRow(_SELECT_ID, id).Scan(&peerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// The peer is not in the database.
			return false
		}
		// Some other error occurred.
		return false
	}

	// If we get here, it means the peer exists in the database.
	return true
}

// Function is equiavalent to ShortString() in libp2p, but it also
// checks if the peer is known in the database and returns the
// node alias if it exists.
// The ShortString() function returns the last 8 chars of the peer ID.
// The input is a full peer ID string.
// Returns the input in case of errors
func getOrCreateNodeAlias(id string) (nodeAlias string) {

	if len(id) < nodeAliasLength {
		return id
	}

	// Nicks to lookup can be less than 8 chars
	if IsKnown(id) {
		nodeAlias, err := LookupNick(id)
		if err == nil {
			return nodeAlias
		}
	}

	return id[len(id)-nodeAliasLength:]

}
