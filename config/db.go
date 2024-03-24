package config

import (
	"github.com/bahner/go-ma-actor/internal"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var DefaultDbPath = internal.NormalisePath(dataHome + "/." + Profile())

func init() {
	pflag.String("db-path", DefaultDbPath, "Directory to use for database.")
	viper.BindPFlag("db.path", pflag.Lookup("db-path"))
	viper.SetDefault("db.path", DefaultDbPath)
}

type DBConfig struct {
	Path string `yaml:"path"`
}

func DB() DBConfig {

	viper.SetDefault("db.path", DefaultDbPath)

	return DBConfig{
		Path: DBPath(),
	}
}

func DBPath() string {
	return viper.GetString("db.path")
}
