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
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	defaultAliasFile   = "~/.ma/aliases.db"
	defaultAliasLength = 8

	SELECT_ENTITY_NICK = "SELECT nick FROM entities WHERE did = ?"
	SELECT_ENTITY_DID  = "SELECT did FROM entities WHERE nick = ?"
	SELECT_NODE_NICK   = "SELECT nick FROM nodes WHERE id = ?"
	SELECT_NODE_ID     = "SELECT id FROM nodes WHERE nick = ?"
	UPSERT_ENTITY      = "INSERT INTO entities (did, nick) VALUES (?, ?) ON CONFLICT(did) DO UPDATE SET nick = ?"
	UPSERT_NODE        = "INSERT INTO nodes (id, nick) VALUES (?, ?) ON CONFLICT(id) DO UPDATE SET nick = ?"
	DELETE_ENTITY      = "DELETE FROM entities WHERE did = ?"
	DELETE_NODE        = "DELETE FROM nodes WHERE id = ?"
)

var (
	once sync.Once
	db   *sql.DB
)

func init() {

	pflag.String("aliases", defaultAliasFile, "File to *write* node aliases to. If the file does not exist, it will be created.")
	viper.BindPFlag("aliases", pflag.Lookup("aliases"))
	viper.SetDefault("aliases", defaultAliasFile)
}

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

	if !did.IsValidDID(id) {
		return "", fmt.Errorf("invalid DID: %s", id)
	}

	db, err := GetDB()
	if err != nil {
		return "", err
	}

	var a string

	err = db.QueryRow(SELECT_ENTITY_NICK, id).Scan(&a)
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

	err = db.QueryRow(SELECT_ENTITY_DID, nick).Scan(&id)
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
		return "", fmt.Errorf("invalid ID: %s", id)
	}

	db, err := GetDB()
	if err != nil {
		return "", err
	}

	var a string

	err = db.QueryRow(SELECT_NODE_NICK, id).Scan(&a)
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

	err = db.QueryRow(SELECT_NODE_ID, nick).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

// Sets an entity alias in the database
// If the alias already exists, it will be updated.
func SetEntityAlias(id string, nick string) error {

	if !did.IsValidDID(id) {
		return fmt.Errorf("invalid DID: %s", id)
	}

	db, err := GetDB()
	if err != nil {
		return err
	}

	_, err = db.Exec(UPSERT_ENTITY, id, nick, nick)
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
		return fmt.Errorf("invalid ID: %s", id)
	}

	db, err := GetDB()
	if err != nil {
		return err
	}

	_, err = db.Exec(UPSERT_NODE, id, nick, nick)
	if err != nil {
		return err
	}

	return nil
}

// Removes an entity alias from the database if it exists
func RemoveEntityAlias(id string) error {

	if !did.IsValidDID(id) {
		return fmt.Errorf("invalid DID: %s", id)
	}

	db, err := GetDB()
	if err != nil {
		return err
	}

	_, err = db.Exec(DELETE_ENTITY, id)
	if err != nil {
		return err
	}

	return nil
}

// Removes a node alias from the database if it exists
func RemoveNodeAlias(id string) error {

	_, err := peer.Decode(id)
	if err != nil {
		return fmt.Errorf("invalid ID: %s", id)
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
		log.Debugf("Error looking up entity alias for %s: %s", id, err)
		return id
	}

	return n
}

// Attempts to look up an entity DID in the database
// If the DID does not exist, the function returns the input string
func LookupEntityNick(nick string) string {

	n, err := GetEntityDID(nick)
	if err != nil {
		log.Debugf("Error looking up DID for entity with nick %s: %s", nick, err)
		return nick
	}

	return n
}

// Attempts to look up a node alias in the database
// If the alias does not exist, the function returns the input string
func LookupNodeID(id string) string {

	// Don't valid ID here. It's not necessary.

	n, err := GetNodeAlias(id)
	if err != nil {
		log.Debugf("Error looking up node alias for NodeID %s: %s", id, err)
		return id
	}

	return n
}

// Attempts to look up a node ID in the database
// If the ID does not exist, the function returns the input string
func LookupNodeNick(nick string) string {

	n, err := GetNodeID(nick)
	if err != nil {
		log.Debugf("Error looking up node ID for nick %s: %s", nick, err)
		return nick
	}

	return n
}

// Returns a string containing all entity aliases
func EntityAliases() string {

	db, err := GetDB()
	if err != nil {
		return fmt.Sprintf("Error getting aliases: %s", err)
	}

	rows, err := db.Query("SELECT did, nick FROM entities")
	if err != nil {
		return fmt.Sprintf("Error getting aliases: %s", err)
	}

	defer rows.Close()

	var s string
	for rows.Next() {
		var did, nick string
		err = rows.Scan(&did, &nick)
		if err != nil {
			return fmt.Sprintf("Error getting aliases: %s", err)
		}
		s += fmt.Sprintf("%s: %s\n", did, nick)
	}

	return s
}

// Returns a string containing all node aliases
func NodeAliases() string {

	db, err := GetDB()
	if err != nil {
		return fmt.Sprintf("Error getting aliases: %s", err)
	}

	rows, err := db.Query("SELECT id, nick FROM nodes")
	if err != nil {
		return fmt.Sprintf("Error getting aliases: %s", err)
	}

	defer rows.Close()

	var s string
	for rows.Next() {
		var id, nick string
		err = rows.Scan(&id, &nick)
		if err != nil {
			return fmt.Sprintf("Error getting aliases: %s", err)
		}
		s += fmt.Sprintf("%s: %s\n", id, nick)
	}

	return s
}

// Creates a node alias from the last 8 characters of the ID
// If the alias already exists, the existing alias is returned
func GetOrCreateNodeAlias(id string) (string, error) {

	_, err := peer.Decode(id)
	if err != nil {
		return "", fmt.Errorf("invalid ID: %s", id)
	}

	// If an alias exists, return it
	n, err := GetNodeAlias(id)

	// If not make one from the last 8 characters of the ID
	// The ID must be more than defaultAliasLength characters long
	if err != nil && len(id) > defaultAliasLength {
		n = id[len(id)-defaultAliasLength:]
		SetNodeAlias(id, n)
		return n, nil
	}

	return n, nil
}
