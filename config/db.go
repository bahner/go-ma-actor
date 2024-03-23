package config

import (
	"github.com/bahner/go-ma-actor/internal"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const defaultDBDirname = ".madb"

var DefaultDbPath = internal.NormalisePath(dataHome + defaultDBDirname)

func init() {
	pflag.String("db-path", DefaultDbPath, "Directory to use for database.")
	viper.BindPFlag("db.path", pflag.Lookup("db-path"))
	viper.SetDefault("db.path", DefaultDbPath)
}

type DBStruct struct {
	Path string `yaml:"path"`
}

type DBConfigStruct struct {
	DB DBStruct `yaml:"db"`
}

func DBConfig() DBConfigStruct {

	viper.SetDefault("db.path", DefaultDbPath)

	return DBConfigStruct{
		DB: DBStruct{
			Path: DBPath(),
		},
	}
}

func DBPath() string {
	return viper.GetString("db.path")
}
