package peer

import (
	"github.com/bahner/go-ma-actor/config/db"
	_ "github.com/mattn/go-sqlite3"
)

const (
	_LOOKUP_ID   = "SELECT id FROM peers WHERE nick = ? or id = ?"
	_LOOKUP_NICK = "SELECT nick FROM peers WHERE nick = ? or id = ?"
	_SELECT_NICK = "SELECT nick FROM peers WHERE id =?"
	_UPSERT_NICK = "INSERT INTO peers (id, nick) VALUES (?, ?) ON CONFLICT(id) DO UPDATE SET nick = ?"
)

const (
	defaultAliasLength = 8
)

// Sets a node in the database
// The key is the node's ID
func SetNickForID(p Peer) error {

	d, err := db.Get()
	if err != nil {
		return err
	}

	_, err = d.Exec(_UPSERT_NICK, p.ID, p.Nick, p.Nick)
	if err != nil {
		return err
	}

	return nil
}

// Fetches a peer nick from the database for the PeerID
func GetNickForID(id string) (string, error) {

	nick := ""

	db, err := db.Get()
	if err != nil {
		return nick, err
	}

	err = db.QueryRow(_SELECT_NICK, id).Scan(&nick)
	if err != nil {
		return nick, err
	}

	return nick, err
}

// Returns the ID of a node by searching for either nick or id
func LookupID(q string) (string, error) {

	id := ""

	db, err := db.Get()
	if err != nil {
		return id, err
	}

	err = db.QueryRow(_LOOKUP_ID, q, q).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil

}

func LookupNick(id string) (string, error) {

	nick := ""

	db, err := db.Get()
	if err != nil {
		return nick, err
	}

	err = db.QueryRow(_LOOKUP_ID, id).Scan(&nick)
	if err != nil {
		return nick, err
	}

	return nick, nil
}

// Looks up a node, but returns the input if it's not found
// So this can be used for non-existing nodes
func Lookup(q string) (string, error) {

	db, err := db.Get()
	if err != nil {
		return q, err
	}

	err = db.QueryRow(_LOOKUP_ID, q, q).Scan(&q)
	if err != nil {
		return q, err
	}

	return q, nil
}
