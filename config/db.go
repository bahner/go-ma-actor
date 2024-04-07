package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/bahner/go-ma-actor/internal"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const CSVMode = 0664

var (
	defaultPeersPath    = internal.NormalisePath(dataHome + "/peers.csv")
	defaultEntitiesPath = internal.NormalisePath(dataHome + "/entities.csv")
	dbFlags             = *pflag.NewFlagSet("db", pflag.ContinueOnError)
	dbOnce              sync.Once
)

func InitDB() {

	dbOnce.Do(func() {

		dbFlags.String("peers", defaultPeersPath, "Filename for CSV peers file.")
		dbFlags.String("entities", defaultEntitiesPath, "Filename for CSV entities file.")
		dbFlags.String("history", defaultHistoryPath(), "Filename for CSV history file.")

		viper.BindPFlag("db.peers", dbFlags.Lookup("peers"))
		viper.BindPFlag("db.entities", dbFlags.Lookup("entities"))
		viper.BindPFlag("db.history", dbFlags.Lookup("history"))

		viper.SetDefault("db.peers", defaultPeersPath)
		viper.SetDefault("db.entities", defaultEntitiesPath)
		viper.SetDefault("db.history", defaultHistoryPath())

		if HelpNeeded() {
			fmt.Println("DB Flags:")
			dbFlags.PrintDefaults()
		} else {
			dbFlags.Parse(os.Args[1:])
		}
	})

}

type DBConfig struct {
	Peers    string `yaml:"peers"`
	Entities string `yaml:"entities"`
	History  string `yaml:"history"`
}

func DB() DBConfig {

	return DBConfig{
		Peers:    DBPeers(),
		Entities: DBEntities(),
		History:  DBHistory(),
	}

}

func DBPeers() string {
	return viper.GetString("db.peers")
}

func DBEntities() string {
	return viper.GetString("db.entities")
}

func DBHistory() string {
	return viper.GetString("db.history")
}

func defaultHistoryPath() string {
	return internal.NormalisePath(dataHome + "/" + Profile() + ".history")
}
