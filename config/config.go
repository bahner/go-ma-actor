package config

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	NAME    = "go-ma-actor"
	VERSION = "v0.0.4"
)

func init() {

	viper.SetConfigName(NAME)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetEnvPrefix(NAME)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %s\n", err)
	}

	pflag.BoolP("version", "v", false, "Print version and exit.")
	viper.BindPFlag("version", pflag.Lookup("version"))

}

// This should be called after pflag.Parse() in main.
func Init() {

	if viper.GetBool("version") {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	InitLogging()
	InitNodeIdentity()
	InitIdentity()
}
