package peer

import (
	"database/sql"
	"errors"

	"github.com/bahner/go-ma-actor/config/db"
	_ "github.com/mattn/go-sqlite3"
)

const (
	_SELECT_ALLOWED = "SELECT allowed FROM peers WHERE id = ?"
	_UPDATE_ALLOWED = "UPDATE peers SET allowed = ? WHERE id = ?"
)

const (
	defaultAllowed = true // We need this to be able to discover newly added peers.
)

// GetAllowedForID returns whether a peer is allowed to be discovered.
// This implies whther the peer is blacklisted or not.
func GetAllowedForID(id string) (bool, error) {

	allowed := bool2int(defaultAllowed)

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

	return int2bool(allowed), nil
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

	_, err = tx.Exec(_UPDATE_ALLOWED, id, bool2int(allowed))
	return err
}

func IsAllowed(id string) bool {
	allowed, err := GetAllowedForID(id)
	if err != nil {
		return defaultAllowed
	}
	return allowed
}

// bool2int converts a bool to an int. true is 1, false is 0.
func bool2int(b bool) int {
	if b {
		return 1
	}
	return 0
}

// Converts an int to a bool. 1 is true, anything else is false.
func int2bool(a int) bool {
	return a == 1
}
