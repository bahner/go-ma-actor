package entity

import (
	"errors"

	"github.com/bahner/go-ma-actor/config/db"
	nanoid "github.com/matoous/go-nanoid/v2"
	_ "github.com/mattn/go-sqlite3"
)

const (
	_SELECT_NICK = "SELECT nick FROM entities WHERE did =?"
	_SELECT_DID  = "SELECT did FROM entities WHERE nick =?"
	_UPSERT      = "INSERT INTO entities (did, nick) VALUES (?, ?) ON CONFLICT(did) DO UPDATE SET nick = ?"
	_DELETE      = "DELETE FROM entities WHERE did = ?"
)

var ErrFailedToSetNick = errors.New("failed to set entity nick")

// Returns the DID . Returns the input if the node does not exist
// This is used before we know in an Entity exists or not. It can be used anywhere.
func GetDID(id string) string {

	did, err := LookupNick(id)
	if err == nil {
		return did
	}

	return id

}

func ListNicks() map[string]string {

	entities := make(map[string]string)

	d, err := db.Get()
	if err != nil {
		return entities
	}

	rows, err := d.Query("SELECT did, nick FROM entities")
	if err != nil {
		return entities
	}
	defer rows.Close()

	for rows.Next() {
		var did, nick string
		err = rows.Scan(&did, &nick)
		if err != nil {
			return entities
		}

		entities[did] = nick
	}

	return entities
}

// Takes a nick as input and returns the corresponding DID
// Else it returns the input as is with an error.
func LookupNick(id string) (string, error) {

	did := id

	d, err := db.Get()
	if err != nil {
		return id, err
	}

	err = d.QueryRow(_SELECT_DID, id).Scan(&did)
	if err != nil {
		return id, err
	}

	return did, nil

}

// Removes a node from the database if it exists
func RemoveNick(id string) error {

	id = GetDID(id) // Check if the reuested ID is a nick

	entities.Delete(id)

	d, err := db.Get()
	if err != nil {
		return err
	}

	_, err = d.Exec(_DELETE, id)
	if err != nil {
		return err
	}

	return nil
}

// Sets a node in the database
// The key is the node's ID
func (e *Entity) SetNick(nick string) error {

	e.Nick = nick

	d, err := db.Get()
	if err != nil {
		return err
	}
	_, err = d.Exec(_UPSERT, e.DID.Id, nick, nick)
	if err != nil {
		return err
	}

	return nil
}

// Fetches a nick from the database or uses the fragment of the DID
func (e *Entity) getOrCreateAndSetNick() error {

	var (
		err error
	)

	// Nicks should be unique. If we have a nick, we're done.
	// Also they must be random initially, so others can't overwrite them.
	db, err := db.Get()
	if err == nil {

		err = db.QueryRow(_SELECT_NICK, e.DID.Id).Scan(&e.Nick)
		if err == nil {
			e.SetNick(e.Nick)
			return nil
		}
	}

	e.Nick, err = nanoid.New()
	if err == nil {
		e.SetNick(e.Nick)
		return nil
	}

	return ErrFailedToSetNick
}
