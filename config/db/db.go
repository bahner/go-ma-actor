package db

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/bahner/go-ma-actor/config"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	dbMaxConnections  = 1 // Required for serialized access to the database
	defaultDBFilename = "ma.db"
	defaultDbTimeout  = 10000
)

var (
	once          sync.Once
	db            *sql.DB
	defaultDbFile = config.DefaultDbFile
)

func init() {

	pflag.String("db-file", defaultDbFile, "File to *write* node peers and entities to. If the file does not exist, it will be created.")
	pflag.Int("db-timeout", defaultDbTimeout, "Timeout for serialized access to the database in milliseconds.")

	viper.BindPFlag("db.file", pflag.Lookup("db-file"))
	viper.BindPFlag("db.timeout", pflag.Lookup("db-timeout"))

}

// Initiates the database connection and creates the tables if they do not exist
func Init() (*sql.DB, error) {

	var onceErr error

	once.Do(func() {

		var err error
		f, err := dbfile()
		if err != nil {
			onceErr = fmt.Errorf("error expanding db file path: %s", err)
			return
		}
		db, err = sql.Open("sqlite3", f)
		if err != nil {
			onceErr = fmt.Errorf("error opening database: %s", err)
			return
		}

		_, err = db.Exec("CREATE TABLE IF NOT EXISTS entities (did VARCHAR(80) PRIMARY KEY, nick VARCHAR(255), UNIQUE(nick) )")
		if err != nil {
			onceErr = fmt.Errorf("error creating entities table: %s", err)
			return
		}

		_, err = db.Exec("CREATE TABLE IF NOT EXISTS peers (id VARCHAR(60) PRIMARY KEY, nick VARCHAR(255), allowed BOOLEAN NOT NULL CHECK (allowed IN (0, 1)), UNIQUE(nick))")
		if err != nil {
			onceErr = fmt.Errorf("error creating peers table: %s", err)
			return
		}

		// Force serialized access to the database with a 100 millisecond timeout. This should be amble time.
		db.SetMaxOpenConns(dbMaxConnections)
		db.Exec("PRAGMA busy_timeout = " + fmt.Sprintf("%d", timeout()))

		log.Infof("Connected to database: %s", f)

	})

	if onceErr != nil {
		return nil, onceErr
	}

	return db, nil
}
func Get() (*sql.DB, error) {

	return Init()
}

// Returns expanded path to the db-file file
// If the expansion fails it returns an empty string
func dbfile() (string, error) {

	p := viper.GetString("db.file")
	p, err := homedir.Expand(p)
	if err != nil {
		return "", err
	}

	return config.NormalisePath(p), nil

}

func timeout() int {
	return viper.GetInt("db.timeout")
}
