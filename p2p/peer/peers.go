package peer

import (
	"sync"

	"github.com/bahner/go-ma-actor/config/db"
	_ "github.com/mattn/go-sqlite3"
)

const (
	_DELETE_PEER = "DELETE FROM peers WHERE id = ?"
	_UPSERT_PEER = "INSERT INTO peers (id, nick, allowed) VALUES (?, ?, ?) ON CONFLICT(id) DO UPDATE SET (nick, allowed) = (?, ?)"
)

// peers safely stores and retrieves Peer values.
var (
	peers sync.Map
)

// Get retrieves a peer's information from the map by ID.
func Get(id string) (Peer, error) {

	var err error

	// Retrieve the peer.
	value, ok := peers.Load(id)
	if ok {
		p := value.(Peer)
		p.Nick, err = GetNickForID(p.ID)
		if err != nil {
			return Peer{}, err
		}
		p.Allowed, err = GetAllowedForID(p.ID)
		if err != nil {
			return Peer{}, err
		}
		return p, nil
	}
	return Peer{}, ErrPeerNotFound
}

// Set modifies an existing peer's information in the map and the database.
func Set(p Peer) error {

	// Delete the peer from the database.
	d, err := db.Get()
	if err != nil {
		return err
	}

	_, err = d.Exec(_UPSERT_PEER, p.ID, p.Nick, p.Allowed, p.Nick, p.Allowed)
	if err != nil {
		return err
	}

	peers.Store(p.ID, p)

	return nil
}

// Delete removes a peer from the map and the database by ID.
func Delete(id string) error {

	// Delete the peer from the database.
	d, err := db.Get()
	if err != nil {
		return err
	}

	_, err = d.Exec(_DELETE_PEER, id)
	if err != nil {
		return err
	}

	// Remove the peer from the sync.Map.
	peers.Delete(id)

	return nil
}

// List returns a slice of all peers in the map.
func List() []Peer {
	p := make([]Peer, 0)
	peers.Range(func(_, value interface{}) bool {
		p = append(p, value.(Peer))
		return true // continue iteration
	})
	return p
}
