package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/bahner/go-ma"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

const (
	NAME       string = "go-ma-actor"
	VERSION    string = "v0.3.1"
	ENV_PREFIX string = "GO_MA_ACTOR"

	configFileMode os.FileMode = 0600
	configDirMode  os.FileMode = 0700
	dataHomeMode   os.FileMode = 0755
)

var (
	configHome        string = xdg.ConfigHome + "/" + ma.NAME + "/"
	dataHome          string = xdg.DataHome + "/" + ma.NAME + "/"
	defaultConfigFile string = NormalisePath(configHome + Profile() + ".yaml")
)

// This should be called after pflag.Parse() in main.
// If you want to use a specific config file, you need to call SetProfile() before Init().
func Init() error {

	var err error

	//VIPER CONFIGURATION

	// Read the config file and environment variables.
	viper.SetEnvPrefix(ENV_PREFIX)
	viper.AutomaticEnv()

	// Handle nested values in environment variables.
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	// Look in the current directory, the home directory and /etc for the config file.
	// In that order.
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(configHome)

	// We *must* read the config file after we have generated the identity.
	// Otherwise: Unforeseen consequences.
	if !generateFlag() {
		log.Infof("Using config file: %s", configFile()) // This one goes to STDERR
		viper.SetConfigFile(configFile())
		err = viper.ReadInConfig()
		if err != nil {
			log.Warnf("No config file found: %s", err)
		}
	}

	// API
	viper.SetDefault("api.maddr", ma.DEFAULT_IPFS_API_MULTIADDR)

	// Logging
	viper.SetDefault("log.file", genDefaultLogFileName(Profile()))
	InitLogging()

	// FLAGS

	// Handle the easy flags first.
	if versionFlag() {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	// Make sure the XDG directories exist before we start writing to them.
	err = createXDGDirectories()
	if err != nil {
		panic(err)
	}

	if generateFlag() {

		// Reinit logging to STDOUT
		log.SetOutput(os.Stdout)
		log.Info("Generating new actor and node identity")
		actor, node := handleGenerateOrExit()
		generateActorConfigFile(actor, node)
		os.Exit(0)
	}

	InitActor()

	// This flag is dependent on the actor to be initialized to make sense.
	if showConfigFlag() {
		Print()
		os.Exit(0)
	}
	return nil

}

func Print() (int, error) {

	configMap := viper.AllSettings()

	configYAML, err := yaml.Marshal(configMap)
	if err != nil {
		panic(err)
	}

	fmt.Println("# " + ActorKeyset().DID.Id)

	return fmt.Println(string(configYAML))
}

func Save() error {

	return viper.WriteConfig()

}

func DataHome() string {
	return dataHome
}

func ConfigHome() string {
	return configHome
}

// Returns the configfile name to use.
// The preferred value is the explcitily requested config file on the command line.
// Else it uses the nick of the actor or the mode.
func configFile() string {

	var (
		filename string
		err      error
	)

	config, err := pflag.CommandLine.GetString("config")
	if err != nil {
		panic(err)
	}

	// Prefer explicitly requested config. If not, use the name of the profile name.
	if config != defaultConfigFile && config != "" {
		filename, err = homedir.Expand(config)
		if err != nil {
			panic(err)
		}
	} else {
		filename = configHome + Profile() + ".yaml"
	}

	return filepath.Clean(filename)

}

func generateFlag() bool {
	// This will exit when done. It will also publish if applicable.
	generateFlag, err := pflag.CommandLine.GetBool("generate")
	if err != nil {
		log.Warnf("config.init: %v", err)
		return false
	}

	return generateFlag
}

func publishFlag() bool {
	publishFlag, err := pflag.CommandLine.GetBool("publish")
	if err != nil {
		log.Warnf("config.init: %v", err)
		return false
	}

	return publishFlag
}

func showConfigFlag() bool {
	showConfigFlag, err := pflag.CommandLine.GetBool("show-config")
	if err != nil {
		log.Warnf("config.init: %v", err)
		return false
	}

	return showConfigFlag
}

func versionFlag() bool {
	versionFlag, err := pflag.CommandLine.GetBool("version")
	if err != nil {
		log.Warnf("config.init: %v", err)
		return false
	}

	return versionFlag
}

func forceFlag() bool {
	forceFlag, err := pflag.CommandLine.GetBool("force")
	if err != nil {
		log.Warnf("config.init: %v", err)
		return false
	}

	return forceFlag
}

func NormalisePath(path string) string {
	return filepath.ToSlash(filepath.Clean(path))
}
