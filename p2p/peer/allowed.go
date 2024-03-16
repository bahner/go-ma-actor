package peer

import (
	"database/sql"
	"errors"

	"github.com/bahner/go-ma-actor/config/db"
	"github.com/bahner/go-ma-actor/internal"
	_ "github.com/mattn/go-sqlite3"
)

const (
	_SELECT_ALLOWED = "SELECT allowed FROM peers WHERE id = ?"
	_UPDATE_ALLOWED = "UPDATE peers SET allowed = ? WHERE id = ?"

	defaultAllowed = true // This is required for discovery to work for hosts that are not in the database.
)

// GetAllowedForID returns whether a peer is allowed to be discovered.
// This implies whther the peer is blacklisted or not.
func GetAllowedForID(id string) (bool, error) {

	allowed := internal.Bool2int(defaultAllowed)

	db, err := db.Get()
	if err != nil {
		return defaultAllowed, err
	}

	err = db.QueryRow(_SELECT_ALLOWED, id).Scan(&allowed)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return defaultAllowed, ErrPeerNotFoundInDB
		}
		return defaultAllowed, err
	}

	return internal.Int2bool(allowed), nil
}

func SetAllowed(id string, allowed bool) error {
	d, err := db.Get()
	if err != nil {
		return err
	}

	tx, err := d.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	_, err = tx.Exec(_UPDATE_ALLOWED, id, internal.Bool2int(allowed))
	return err
}

func IsAllowed(id string) bool {
	allowed, err := GetAllowedForID(id)
	if err != nil {
		return defaultAllowed
	}
	return allowed
}
