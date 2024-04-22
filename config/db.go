package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/bahner/go-ma-actor/internal"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	CSVMode = 0664
)

var (
	err    error
	dbOnce sync.Once

	dbFlags = *pflag.NewFlagSet("db", pflag.ContinueOnError)
)

func InitDB() {

	dbOnce.Do(func() {

		dbFlags.String("entities", defaultEntitiesPath(), "Filename for CSV entities file.")
		dbFlags.String("history", defaultHistoryPath(), "Filename for CSV history file.")
		dbFlags.String("keystore", defaultKeystorePath(), "Folder name to store keys in.")
		dbFlags.String("peers", defaultPeersPath(), "Filename for CSV peers file.")

		viper.BindPFlag("db.entities", dbFlags.Lookup("entities"))
		viper.BindPFlag("db.history", dbFlags.Lookup("history"))
		viper.BindPFlag("db.keystore", dbFlags.Lookup("keystore"))
		viper.BindPFlag("db.peers", dbFlags.Lookup("peers"))

		viper.SetDefault("db.entities", defaultEntitiesPath())
		viper.SetDefault("db.history", defaultHistoryPath())
		viper.SetDefault("db.keystore", defaultKeystorePath())
		viper.SetDefault("db.peers", defaultPeersPath())

		if HelpNeeded() {
			fmt.Println("DB Flags:")
			dbFlags.PrintDefaults()
		} else {
			dbFlags.Parse(os.Args[1:])
		}

	})

}

type DBConfig struct {
	Entities string `yaml:"entities"`
	History  string `yaml:"history"`
	Keystore string `yaml:"keystore"`
	Peers    string `yaml:"peers"`
}

func DB() DBConfig {

	return DBConfig{
		Entities: DBEntities(),
		History:  DBHistory(),
		Keystore: DBKeystore(),
		Peers:    DBPeers(),
	}

}

func DBEntities() string {
	return viper.GetString("db.entities")
}

func DBHistory() string {
	return viper.GetString("db.history")
}

func DBKeystore() string {
	return viper.GetString("db.keystore")
}

func DBPeers() string {
	return viper.GetString("db.peers")
}

func defaultHistoryPath() string {
	return internal.NormalisePath(dataHome + "/" + Profile() + ".history")
}

func defaultEntitiesPath() string {
	return internal.NormalisePath(dataHome + "/entities.csv")
}

func defaultPeersPath() string {
	return internal.NormalisePath(dataHome + "/peers.csv")
}
