package db

import (
	"fmt"
	"sync"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/internal"
	badger "github.com/dgraph-io/badger/v3"
	"github.com/mitchellh/go-homedir"
)

var (
	db   *badger.DB
	once sync.Once
)

// InitDB initializes the Badger database and sets the global `db` variable.
// It uses `sync.Once` to ensure that the database is opened only once.
func initDB() (*badger.DB, error) {

	var err error

	once.Do(func() {

		var dbPath string

		dbPath, err = dbDir()
		if err != nil {
			return
		}

		read_only := config.DBRo()

		badgerOpts := badger.DefaultOptions(dbPath).WithReadOnly(read_only)

		db, err = badger.Open(badgerOpts)
		if err != nil {
			return
		}
	})

	return db, nil
}

// CloseDB closes the Badger database. This should be called when your application exits.
func Close() (err error) {
	if db != nil {
		err = db.Close()
		if err != nil {
			return fmt.Errorf("failed to close BadgerDB: %v", err)
		}
	}
	return nil
}

func DB() *badger.DB {

	if db == nil {
		initDB()
	}

	return db
}

// Returns expanded path to the dbDir file
// If the expansion fails it returns an empty string
func dbDir() (string, error) {

	p := config.DBPath()
	p, err := homedir.Expand(p)
	if err != nil {
		return "", err
	}

	return internal.NormalisePath(p), nil

}
