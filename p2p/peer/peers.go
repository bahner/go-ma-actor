package peer

import (
	"sync"

	"github.com/bahner/go-ma-actor/config/db"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

const (
	_DELETE_PEER = "DELETE FROM peers WHERE id = ?"
	_INSERT_PEER = "INSERT INTO peers (id, nick, allowed) VALUES (?, ?, ?)"
	// NB! UPSERT doesn't work properly with TEXT columns in sqlite3 it seems
	// SO we have to delete and insert instead of update
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

	_, err = tx.Exec(_DELETE_PEER, p.ID)
	if err != nil {
		log.Debugf("Failed to set peer %s: %s", p.ID, err.Error())
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(_INSERT_PEER, p.ID, p.Nick, sqlAllowed)
	if err != nil {
		log.Debugf("Failed to set peer %s: %s", p.ID, err.Error())
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

// CachedPeers returns a slice of all peers in the the DB.
// NB! This is not the database and should be used to check nick and allowed status.
// It's just a cache.
func CachedPeers() ([]Peer, error) {
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
