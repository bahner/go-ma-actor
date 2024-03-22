package db

import (
	"fmt"
	"sync"

	"github.com/bahner/go-ma-actor/config"
	badger "github.com/dgraph-io/badger/v3"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {

	pflag.String("db-path", config.DefaultDbPath, "Directory to use for database.")
	viper.BindPFlag("db.path", pflag.Lookup("db-path"))
	viper.SetDefault("db.path", config.DefaultDbPath)

}

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

		db, err = badger.Open(badger.DefaultOptions(dbPath))
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

	p := viper.GetString("db.path")
	p, err := homedir.Expand(p)
	if err != nil {
		return "", err
	}

	return config.NormalisePath(p), nil

}
