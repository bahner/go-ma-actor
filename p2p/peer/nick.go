package peer

import (
	"database/sql"
	"errors"

	"github.com/bahner/go-ma-actor/config/db"
	_ "github.com/mattn/go-sqlite3"
)

const (
	_LOOKUP_NICK  = "SELECT nick FROM peers WHERE nick = ? OR id = ?"
	_SELECT_NICK  = "SELECT nick FROM peers WHERE id = ?"
	_UPSERT_NICK  = "INSERT INTO peers (id, nick) VALUES (?, ?) ON CONFLICT(id) DO UPDATE SET nick = excluded.nick;"
	_SELECT_NICKS = "SELECT id, nick FROM peers"
)

var nodeAliasLength = 8

// SetNickForID updates or inserts the nick for a given peer ID, using a transaction.
func SetNickForID(id string, nick string) error {

	var err error

	d, err := db.Get()
	if err != nil {
		return err
	}

	tx, err := d.Begin()
	if err != nil {
		return err
	}

	// Defer a function to ensure the transaction is either committed or rolled back.
	defer func() {
		// Only attempt to rollback or commit if `err` is not already set.
		if err != nil {
			tx.Rollback()
		} else {
			commitErr := tx.Commit()
			if commitErr != nil {
				err = commitErr
			}
		}
	}()

	_, err = tx.Exec(_UPSERT_NICK, id, nick)
	return err
}

// GetNickForID retrieves a peer's nickname by its ID.
func GetNickForID(id string) (string, error) {
	db, err := db.Get()
	if err != nil {
		return "", err
	}

	var nick string
	err = db.QueryRow(_SELECT_NICK, id).Scan(&nick)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrPeerNotFoundInDB
		}
		return "", err
	}

	return nick, nil
}

// LookupNick finds a nickname for a peerby its ID or Nick.
func LookupNick(id string) (string, error) {
	db, err := db.Get()
	if err != nil {
		return "", err
	}

	var nick string
	err = db.QueryRow(_LOOKUP_NICK, id, id).Scan(&nick)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrPeerNotFoundInDB
		}
		return "", err
	}

	return nick, nil
}
func Nicks() map[string]string {
	db, err := db.Get()
	if err != nil {
		return nil
	}

	rows, err := db.Query(_SELECT_NICKS)
	if err != nil {
		return nil
	}
	defer rows.Close()

	peers := make(map[string]string)
	for rows.Next() {
		var id, nick string
		err = rows.Scan(&id, &nick)
		if err != nil {
			return peers
		}

		peers[id] = nick
	}

	return peers
}

// Function is equiavalent to ShortString() in libp2p, but it also
// checks if the peer is known in the database and returns the
// node alias if it exists.
// The ShortString() function returns the last 8 chars of the peer ID.
// The input is a full peer ID string.
// Returns the input in case of errors
func GetOrCreateNick(id string) (nodeAlias string) {

	// If we find the ID, fetch the nick and return it.
	id, err := LookupID(id)
	if err == nil {
		nodeAlias, err := LookupNick(id)
		if err == nil && nodeAlias != "" {
			return nodeAlias
		}
	}

	// We can't shorten the ID if it's too short.
	if len(id) < nodeAliasLength {
		return id
	}

	return id[len(id)-nodeAliasLength:]

}
