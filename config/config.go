package config

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	NAME       = "go-ma-actor"
	VERSION    = "v0.0.4"
	ENV_PREFIX = "GO_MA_ACTOR"
)

var configFile string = ""

func init() {

	// Look in the current directory, the home directory and /etc for the config file.
	// In that order.
	viper.SetConfigName(NAME)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.ma")
	viper.AddConfigPath("/etc/ma")

	viper.SetEnvPrefix(ENV_PREFIX)
	viper.AutomaticEnv()

	// Allow to set config file via command line flag.
	pflag.StringVarP(&configFile, "config", "c", "", "Config file to use.")

	pflag.BoolP("version", "v", false, "Print version and exit.")
	viper.BindPFlag("version", pflag.Lookup("version"))

	pflag.String("loglevel", defaultLogLevel, "Loglevel to use for application.")
	viper.BindPFlag("log.level", pflag.Lookup("loglevel"))

	pflag.String("logfile", defaultLogfile, "Logfile to use for application.")
	viper.BindPFlag("log.file", pflag.Lookup("logfile"))

}

// This should be called after pflag.Parse() in main.
// The name parameter is the name of the config file to search for.
func Init(config string) error {

	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else if config != "" {
		viper.SetConfigFile(config)
	}

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("fatal error config file: %s", err)
	}

	if viper.GetBool("version") {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	return nil

}
