package peer

import (
	"database/sql"
	"errors"

	"github.com/bahner/go-ma-actor/config/db"
	_ "github.com/mattn/go-sqlite3"
)

const (
	_LOOKUP_ID    = "SELECT id FROM peers WHERE nick = ? OR id = ?"
	_LOOKUP_NICK  = "SELECT nick FROM peers WHERE nick = ? OR id = ?"
	_SELECT_ID    = "SELECT id FROM peers WHERE nick = ?"
	_SELECT_NICK  = "SELECT nick FROM peers WHERE id = ?"
	_UPDATE_NICK  = "UPDATE peers SET nick = ? WHERE id = ?"
	_SELECT_NICKS = "SELECT id, nick FROM peers"
)

var (
	ErrPeerNotFoundInDB    = errors.New("peer not found in database")
	ErrDBTransactionFailed = errors.New("database transaction failed")
)

// SetNickForID updates or inserts the nick for a given peer ID, using a transaction.
func SetNickForID(p Peer) error {

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

	_, err = tx.Exec(_UPDATE_NICK, p.ID, p.Nick)
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

// Lookup finds a peer nickname by its ID or Nick.
// If the name is not found, it returns the input name.
func Lookup(name string) (string, error) {

	id, err := LookupID(name)
	if err != nil {
		return name, ErrPeerNotFoundInDB
	}

	return id, nil
}

// Return a boolean whther the peer is known not
// This this should err on the side of caution and return false
// The input can be a peer ID or a nickname.
func IsKnown(id string) bool {
	db, err := db.Get()
	if err != nil {
		return false
	}

	// We just need to know if the peer exists, so we select the id itself.
	var peerID string
	err = db.QueryRow(_LOOKUP_ID, id).Scan(&peerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false
		}
		// Some other error occurred.
		return false
	}

	// If we get here, it means the peer exists in the database.
	return true
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
func getOrCreateNick(id string) (nodeAlias string) {

	// If we find the ID, fetch the nick and return it.
	id, err := LookupID(id)
	if err == nil {
		nodeAlias, err := LookupNick(id)
		if err == nil {
			return nodeAlias
		}
	}

	// Else return the last 8 chars of the peer ID.
	return id[len(id)-nodeAliasLength:]

}
