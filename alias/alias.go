package alias

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma/did"
	"github.com/libp2p/go-libp2p/core/peer"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

const (
	defaultAliasLength = 8

	_SELECT_ENTITY_NICK = "SELECT nick FROM entities WHERE did = ? or nick = ?"
	_SELECT_ENTITY_DID  = "SELECT did FROM entities WHERE did = ? or nick = ?"
	_SELECT_NODE_NICK   = "SELECT nick FROM nodes WHERE id = ? or nick = ?"
	_SELECT_NODE_ID     = "SELECT id FROM nodes WHERE id =? or nick = ?"
	_UPSERT_ENTITY      = "INSERT INTO entities (did, nick) VALUES (?, ?) ON CONFLICT(did) DO UPDATE SET nick = ?"
	_UPSERT_NODE        = "INSERT INTO nodes (id, nick) VALUES (?, ?) ON CONFLICT(id) DO UPDATE SET nick = ?"
	_DELETE_ENTITY      = "DELETE FROM entities WHERE did = ? or nick = ?"
	_DELETE_NODE        = "DELETE FROM nodes WHERE id = ? or nick = ?"
)

var (
	once sync.Once
	db   *sql.DB
)

// Initiates the database connection and creates the tables if they do not exist
func GetDB() (*sql.DB, error) {

	var onceErr error

	once.Do(func() {
		var err error
		db, err = sql.Open("sqlite3", config.GetAliases())
		if err != nil {
			onceErr = fmt.Errorf("error opening database: %s", err)
			return
		}

		_, err = db.Exec("CREATE TABLE IF NOT EXISTS entities (did TEXT PRIMARY KEY, nick TEXT)")
		if err != nil {
			onceErr = fmt.Errorf("error creating entities table: %s", err)
			return
		}

		_, err = db.Exec("CREATE TABLE IF NOT EXISTS nodes (id TEXT PRIMARY KEY, nick TEXT)")
		if err != nil {
			onceErr = fmt.Errorf("error creating nodes table: %s", err)
			return
		}

	})

	if onceErr != nil {
		return nil, onceErr
	}

	return db, nil
}

// Fetches an entity alias from the database
// Returns an empty string an error if the alias does not exist
func GetEntityAlias(id string) (string, error) {

	err := did.Validate(id)
	if err != nil {
		return "", fmt.Errorf("GetEntityAlias %s: %w", id, err)
	}

	db, err := GetDB()
	if err != nil {
		return "", err
	}

	var a string

	err = db.QueryRow(_SELECT_ENTITY_NICK, id, id).Scan(&a)
	if err != nil {
		return "", err
	}

	return a, nil
}

// Fetches a n entity DID from the database
// Returns an empty string an error if the alias does not exist
func GetEntityDID(nick string) (string, error) {

	if nick == "" {
		return "", fmt.Errorf("query is empty: %s", nick)
	}

	db, err := GetDB()
	if err != nil {
		return "", err
	}

	var id string

	err = db.QueryRow(_SELECT_ENTITY_DID, nick, nick).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

// Fetches a node alias from the database
// Returns an empty string an error if the alias does not exist
func GetNodeAlias(id string) (string, error) {

	_, err := peer.Decode(id)
	if err != nil {
		return "", fmt.Errorf("GetNodeAlias %s, %w", id, err)
	}

	db, err := GetDB()
	if err != nil {
		return "", err
	}

	var a string

	err = db.QueryRow(_SELECT_NODE_NICK, id, id).Scan(&a)
	if err != nil {
		return "", err
	}

	return a, nil
}

// Fetches the ID of a node from the database
// Returns an empty string an error if the alias does not exist
func GetNodeID(nick string) (string, error) {

	if nick == "" {
		return "", fmt.Errorf("query is empty: %s", nick)
	}

	db, err := GetDB()
	if err != nil {
		return "", err
	}

	var id string

	err = db.QueryRow(_SELECT_NODE_ID, nick, nick).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

// Sets an entity alias in the database
// If the alias already exists, it will be updated.
func SetEntityAlias(id string, nick string) error {

	err := did.Validate(id)
	if err != nil {
		return fmt.Errorf("SetEntityAlias %s: %w", id, err)
	}

	db, err := GetDB()
	if err != nil {
		return err
	}

	_, err = db.Exec(_UPSERT_ENTITY, id, nick, nick)
	if err != nil {
		return err
	}

	return nil
}

// Sets a node alias in the database
// If the alias already exists, it will be updated.
func SetNodeAlias(id string, nick string) error {

	_, err := peer.Decode(id)
	if err != nil {
		return fmt.Errorf("invalid Peer ID %s: %w", id, err)
	}

	db, err := GetDB()
	if err != nil {
		return err
	}

	_, err = db.Exec(_UPSERT_NODE, id, nick, nick)
	if err != nil {
		return err
	}

	return nil
}

// Removes an entity alias from the database if it exists
func RemoveEntityAlias(id string) error {

	err := did.Validate(id)
	if err != nil {
		return fmt.Errorf("RemoveEntityAlias %s: %w", id, err)
	}

	db, err := GetDB()
	if err != nil {
		return err
	}

	_, err = db.Exec(_DELETE_ENTITY, id)
	if err != nil {
		return err
	}

	return nil
}

// Removes a node alias from the database if it exists
func RemoveNodeAlias(id string) error {

	_, err := peer.Decode(id)
	if err != nil {
		return fmt.Errorf("invalid ID %s: %w", id, err)
	}

	db, err := GetDB()
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM nodes WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}

// Attempts to look up an entity alias in the database
// If the alias does not exist, the function returns the input string
func LookupEntityDID(id string) string {

	n, err := GetEntityAlias(id)
	if err != nil {
		log.Debugf("Error looking up entity DID for %s: %v", id, err)
		return id
	}

	return n
}

// Attempts to look up an entity DID in the database
// If the DID does not exist, the function returns the input string
func LookupEntityNick(nick string) string {

	n, err := GetEntityDID(nick)
	if err != nil {
		log.Debugf("Error looking up DID for entity with nick %s: %v", nick, err)
		return nick
	}

	return n
}

// Attempts to look up a node alias in the database
// If the alias does not exist, the function returns the input string
func LookupNodeID(id string) string {

	// Don't validate ID here. It's not necessary.

	n, err := GetNodeAlias(id)
	if err != nil {
		log.Debugf("Error looking up node alias for NodeID %s: %v", id, err)
		return id
	}

	return n
}

// Attempts to look up a node ID in the database
// If the ID does not exist, the function returns the input string
func LookupNodeAlias(nick string) string {

	n, err := GetNodeID(nick)
	if err != nil {
		log.Debugf("Error looking up node ID for nick %s: %v", nick, err)
		return nick
	}

	return n
}

// Returns a string containing all entity aliases
func EntityAliases() string {

	db, err := GetDB()
	if err != nil {
		return fmt.Sprintf("Error getting aliases: %v", err)
	}

	rows, err := db.Query("SELECT did, nick FROM entities")
	if err != nil {
		return fmt.Sprintf("Error getting aliases: %v", err)
	}

	defer rows.Close()

	var s string
	for rows.Next() {
		var did, nick string
		err = rows.Scan(&did, &nick)
		if err != nil {
			return fmt.Sprintf("Error getting aliases: %v", err)
		}
		s += fmt.Sprintf("%s: %s\n", did, nick)
	}

	return s
}

// Returns a string containing all node aliases
func NodeAliases() string {

	db, err := GetDB()
	if err != nil {
		return fmt.Sprintf("Error getting aliases: %v", err)
	}

	rows, err := db.Query("SELECT id, nick FROM nodes")
	if err != nil {
		return fmt.Sprintf("Error getting aliases: %v", err)
	}

	defer rows.Close()

	var s string
	for rows.Next() {
		var id, nick string
		err = rows.Scan(&id, &nick)
		if err != nil {
			return fmt.Sprintf("Error getting aliases: %v", err)
		}
		s += fmt.Sprintf("%s: %s\n", id, nick)
	}

	return s
}

// Creates a node alias from the last 8 characters of the ID
// If the alias already exists, the existing alias is returned
// Case of error the input is returned
func GetOrCreateNodeAlias(id string) string {

	_, err := peer.Decode(id)
	if err != nil {
		return id
	}

	// If an alias exists, return it
	n, err := GetNodeAlias(id)

	// If not make one from the last 8 characters of the ID
	// The ID must be more than defaultAliasLength characters long
	if err != nil && len(id) > defaultAliasLength {
		n = id[len(id)-defaultAliasLength:]
		SetNodeAlias(id, n)
		return n
	}

	return n
}

// Creates a node alias from the last 8 characters of the ID
// If the alias already exists, the existing alias is returned
// If the id is an existing alias, the alias is returned
func GetOrCreateEntityAlias(id string) string {

	// If the DID is valid, use it for a lookup or alias creation
	d, err := did.New(id)
	if err == nil {
		// Lookup any existing alias
		n, err := GetEntityAlias(d.Id)
		if err == nil && n != "" {
			return n
		}

		// Use the fragment as the alias
		err = SetEntityAlias(id, d.Fragment)
		if err == nil {
			return d.Fragment
		}

	}

	// Return the input, if all else fails
	return id
}
