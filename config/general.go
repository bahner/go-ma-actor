package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	NAME    = "go-ma-actor"
	VERSION = "v0.0.4"
)

func init() {
	pflag.BoolP("version", "v", false, "Print version and exit.")
	viper.BindPFlag("version", pflag.Lookup("version"))
}

func GetName() string {
	return NAME
}

func GetVersion() string {
	return VERSION
}
