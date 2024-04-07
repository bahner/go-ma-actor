package config

import (
	"github.com/bahner/go-ma-actor/internal"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const CSVMode = 0664

var (
	defaultPeersPath    = internal.NormalisePath(dataHome + "/peers.csv")
	defaultEntitiesPath = internal.NormalisePath(dataHome + "/entities.csv")
	defaultHistoryPath  = internal.NormalisePath(dataHome + "/" + Profile() + ".history")
)

func init() {
	pflag.String("peers", defaultPeersPath, "Filename for CSV peers file.")
	pflag.String("entities", defaultEntitiesPath, "Filename for CSV entities file.")
	pflag.String("history", defaultHistoryPath, "Filename for CSV history file.")

	viper.BindPFlag("db.peers", pflag.Lookup("peers"))
	viper.BindPFlag("db.entities", pflag.Lookup("entities"))
	viper.BindPFlag("db.history", pflag.Lookup("history"))

	viper.SetDefault("db.peers", defaultPeersPath)
	viper.SetDefault("db.entities", defaultEntitiesPath)
	viper.SetDefault("db.history", defaultHistoryPath)
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
