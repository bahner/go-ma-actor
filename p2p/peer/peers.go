package peer

import (
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

// Set modifies an existing peer's information in the map and the database.
func Set(id string, nick string, allowed bool) error {
	d, err := db.Get()
	if err != nil {
		return err
	}

	tx, err := d.Begin()
	if err != nil {
		return err
	}

	sqlAllowed := bool2int(allowed)

	_, err = tx.Exec(_DELETE_PEER, id)
	if err != nil {
		log.Debugf("Failed to set peer %s: %s", id, err.Error())
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(_INSERT_PEER, id, nick, sqlAllowed)
	if err != nil {
		log.Debugf("Failed to set peer %s: %s", id, err.Error())
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

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

	return nil
}
