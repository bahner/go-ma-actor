package peer

import (
	"database/sql"
	"errors"

	"github.com/bahner/go-ma-actor/config/db"
	_ "github.com/mattn/go-sqlite3"
)

const (
	defaultAliasLength = 8

	_LOOKUP_ID   = "SELECT id FROM peers WHERE nick = ? OR id = ?"
	_LOOKUP_NICK = "SELECT nick FROM peers WHERE nick = ? OR id = ?"
	_SELECT_ID   = "SELECT id FROM peers WHERE nick = ?"
	_SELECT_NICK = "SELECT nick FROM peers WHERE id = ?"
	_UPSERT_NICK = "INSERT INTO peers (id, nick) VALUES (?, ?) ON CONFLICT(id) DO UPDATE SET nick = EXCLUDED.nick"
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

	_, err = tx.Exec(_UPSERT_NICK, p.ID, p.Nick, p.Nick)
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

// LookupNick finds a peer nickname by its ID or Nick.
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
