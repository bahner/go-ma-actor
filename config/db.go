package config

import "github.com/spf13/viper"

const defaultDBDirname = ".madb"

var DefaultDbPath = NormalisePath(dataHome + defaultDBDirname)

func DBPath() string {
	return viper.GetString("db.path")
}
