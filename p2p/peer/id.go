package peer

import (
	"database/sql"
	"errors"

	"github.com/bahner/go-ma-actor/config/db"
	_ "github.com/mattn/go-sqlite3"
)

const (
	_LOOKUP_ID = "SELECT id FROM peers WHERE nick = ? OR id = ?"
	_SELECT_ID = "SELECT id FROM peers WHERE nick = ?"
)

var (
	ErrPeerNotFoundInDB    = errors.New("peer not found in database")
	ErrDBTransactionFailed = errors.New("database transaction failed")
)

// GetIDForNick retrieves a peer's ID by its nickname.
func GetIDForNick(nick string) (string, error) {
	db, err := db.Get()
	if err != nil {
		return "", err
	}

	var id string
	err = db.QueryRow(_SELECT_ID, nick).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrPeerNotFoundInDB
		}
		return "", err
	}

	return id, nil
}

// LookupID finds a peer ID by its nickname or ID.
// Returns the input string if the peer is not found.
func LookupID(q string) (string, error) {
	db, err := db.Get()
	if err != nil {
		return "", err
	}

	var id string
	err = db.QueryRow(_LOOKUP_ID, q, q).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return q, ErrPeerNotFoundInDB
		}
		return q, err
	}

	return id, nil
}

// Lookup finds a peer nickname by its ID or Nick.
// If the name is not found, it returns the input name.
func Lookup(name string) string {

	id, err := LookupID(name)
	if err != nil {
		return name
	}

	return id
}

// Return a boolean whther the peer is known not
// This this should err on the side of caution and return false
// The input can be a peer ID or a nickname.
func IsKnown(id string) bool {

	_, err := LookupID(id)

	return err == nil
}
