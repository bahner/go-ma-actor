package config

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

const (
	NAME       = "go-ma-actor"
	VERSION    = "v0.2.1"
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
	pflag.Bool("show-config", false, "Whether to print the config.")
	viper.BindPFlag("show-config", pflag.Lookup("show-config"))
	pflag.Bool("show-defaults", false, "Whether to print the config.")
	viper.BindPFlag("show-defaults", pflag.Lookup("show-defaults"))

	pflag.BoolP("version", "v", false, "Print version and exit.")
	viper.BindPFlag("version", pflag.Lookup("version"))

}

// This should be called after pflag.Parse() in main.
// The name parameter is the name of the config file to search for without the extension.
func Init(configName string) error {

	if configFile != "" {
		configFile, err := homedir.Expand(configFile)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		log.Infof("Using config file: %s", configFile)
		viper.SetConfigFile(configFile)
	} else if configName != "" {
		viper.SetConfigName(configName)
	}

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("No config file found: %s", err)
	}

	if viper.GetBool("version") {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	// This will exit when done. It will also publish if applicable.
	if viper.GetBool("generate") {
		log.Info("Generating new keyset and node identity")
		actor, node := handleGenerateOrExit()
		generateConfigFile(actor, node)
		os.Exit(0)
	}

	if viper.GetBool("show-config") {
		configMap := viper.AllSettings()
		configYAML, err := yaml.Marshal(configMap)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		// Print the YAML to stdout or write it to a template file
		fmt.Println(string(configYAML))
		os.Exit(0)
	}

	if viper.GetBool("show-defaults") {

		// Print the YAML to stdout or write it to a template file
		generateConfigFile("zNO_DEFAULT_ACTOR_IDENITY", "zNO_DEFAULT_NODE_IDENITY")
		os.Exit(0)
	}

	InitLogging()
	InitP2P()

	if !RelayMode() {
		InitActor()
	}

	return nil

}
