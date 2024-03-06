package peer

import (
	"sync"

	"github.com/bahner/go-ma-actor/config/db"
	_ "github.com/mattn/go-sqlite3"
)

const (
	_DELETE_PEER = "DELETE FROM peers WHERE id = ?"
	_UPSERT_PEER = "INSERT INTO peers (id, nick, allowed) VALUES (?, ?, ?) ON CONFLICT(id) DO UPDATE SET nick = EXCLUDED.nick, allowed = EXCLUDED.allowed"
)

var peers sync.Map

// Get retrieves a peer's information from the map by ID.
func Get(id string) (Peer, error) {
	value, ok := peers.Load(id)
	if !ok {
		return Peer{}, ErrPeerNotFound
	}
	p, ok := value.(Peer)
	if !ok {
		// This should not happen if all stored values are of type Peer
		return Peer{}, ErrInvalidPeerType
	}

	var err error
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

// Set modifies an existing peer's information in the map and the database.
func Set(p Peer) error {
	d, err := db.Get()
	if err != nil {
		return err
	}

	tx, err := d.Begin()
	if err != nil {
		return err
	}

	sqlAllowed := bool2int(p.Allowed)

	_, err = tx.Exec(_UPSERT_PEER, p.ID, p.Nick, sqlAllowed)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	peers.Store(p.ID, p)
	return nil
}

// Delete removes a peer from the map and the database by ID.
func Delete(id string) error {
	d, err := db.Get()
	if err != nil {
		return err
	}

	tx, err := d.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(_DELETE_PEER, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	peers.Delete(id)
	return nil
}

// List returns a slice of all peers in the map.
func List() ([]Peer, error) {
	var pList []Peer
	peers.Range(func(_, value interface{}) bool {
		p, ok := value.(Peer)
		if !ok {
			// This should not happen if all stored values are of type Peer
			return false // stop iteration
		}
		pList = append(pList, p)
		return true // continue iteration
	})
	return pList, nil
}
