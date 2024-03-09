package peer

import (
	"database/sql"
	"errors"

	"github.com/bahner/go-ma-actor/config/db"
	_ "github.com/mattn/go-sqlite3"
)

const (
	_LOOKUP_ID  = "SELECT id FROM peers WHERE nick = ? OR id = ?"
	_SELECT_ID  = "SELECT id FROM peers WHERE nick = ?"
	_SELECT_IDS = "SELECT id FROM peers"
)

var (
	ErrPeerNotFoundInDB    = errors.New("peer not found in database")
	ErrDBTransactionFailed = errors.New("database transaction failed")
)

// LookupID finds a peer ID by its nickname or ID.
func LookupID(q string) (string, error) {
	db, err := db.Get()
	if err != nil {
		return "", err
	}

	var id string
	err = db.QueryRow(_LOOKUP_ID, q, q).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrPeerNotFoundInDB
		}
		return "", err
	}

	return id, nil
}

// Returns a slic of all known peer IDs.
func IDS() ([]string, error) {
	db, err := db.Get()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(_SELECT_IDS)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
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
