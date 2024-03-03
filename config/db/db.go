package db

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/bahner/go-ma-actor/config"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	dbMaxConnections  = 1 // Required for serialized access to the database
	defaultDBFilename = "ma.db"
	defaultDbTimeout  = 100
)

var (
	once          sync.Once
	db            *sql.DB
	defaultDbFile = config.DefaultDBFile()
)

func init() {

	pflag.String("db-file", defaultDbFile, "File to *write* node peers and entities to. If the file does not exist, it will be created.")
	viper.BindPFlag("db.file", pflag.Lookup("db-file"))
	viper.SetDefault("db.file", defaultDbFile)

	pflag.Int("db-timeout", defaultDbTimeout, "Timeout for serialized access to the database in milliseconds.")
	viper.BindPFlag("db.timeout", pflag.Lookup("db-timeout"))
	viper.SetDefault("db.timeout", defaultDbTimeout)
}

// Initiates the database connection and creates the tables if they do not exist
func Get() (*sql.DB, error) {

	var onceErr error

	once.Do(func() {
		var err error
		db, err = sql.Open("sqlite3", dbfile())
		if err != nil {
			onceErr = fmt.Errorf("error opening database: %s", err)
			return
		}

		_, err = db.Exec("CREATE TABLE IF NOT EXISTS entities (did VARCHAR(80) PRIMARY KEY, nick VARCHAR(255), UNIQUE(nick) )")
		if err != nil {
			onceErr = fmt.Errorf("error creating entities table: %s", err)
			return
		}

		_, err = db.Exec("CREATE TABLE IF NOT EXISTS nodes (id VARCHAR(60) PRIMARY KEY, node BLOB NO NULL)")
		if err != nil {
			onceErr = fmt.Errorf("error creating nodes table: %s", err)
			return
		}

		// Force serialized access to the database with a 100 millisecond timeout. This should be amble time.
		db.SetMaxOpenConns(dbMaxConnections)
		db.Exec("PRAGMA busy_timeout = " + fmt.Sprintf("%d", timeout()))

	})

	if onceErr != nil {
		return nil, onceErr
	}

	return db, nil
}

// Returns expanded path to the db-file file
// If the expansion fails it returns an empty string
func dbfile() string {

	path := viper.GetString("db.file")
	path, err := homedir.Expand(path)
	if err != nil {
		return ""
	}

	return filepath.FromSlash(filepath.Clean(path))

}

func timeout() int {
	return viper.GetInt("db.timeout")
}
