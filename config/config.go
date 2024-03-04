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
	VERSION    string = "v0.2.4"
	ENV_PREFIX string = "GO_MA_ACTOR"

	configFileMode os.FileMode = 0600
	configDirMode  os.FileMode = 0700
	dataHomeMode   os.FileMode = 0755
)

var (
	config            string = ""
	configHome        string = xdg.ConfigHome + "/" + ma.NAME + "/"
	dataHome          string = xdg.DataHome + "/" + ma.NAME + "/"
	defaultConfigFile string = NormalisePath(configHome + defaultActor + ".yaml")
)

func init() {

	// Allow to set config file via command line flag.
	pflag.StringVarP(&config, "config", "c", defaultConfigFile, "Config file to use.")

	pflag.Bool("force", false, "Whether to force any operation, eg. file overwrite")
	pflag.Bool("show-config", false, "Whether to print the config.")
	pflag.Bool("show-defaults", false, "Whether to print the config.")
	pflag.BoolP("version", "v", false, "Print version and exit.")

}

// This should be called after pflag.Parse() in main.
// The name parameter is the name of the config file to search for without the extension.
func Init(mode string) error {

	// Read the config file and environment variables.
	viper.SetEnvPrefix(ENV_PREFIX)
	viper.AutomaticEnv()

	// Handle nested values in environment variables.
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	// We need the nick to be set before we can generate the config file correctly.
	viper.BindPFlag("actor.nick", pflag.Lookup("nick"))

	// Handle the easy flags first.
	if versionFlag() {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	if showDefaultsFlag() {
		// Print the YAML to stdout or write it to a template file
		generateActorConfigFile(fakeActorIdentity, fakeP2PIdentity)
		os.Exit(0)
	}

	// Now we need logging to be set up, so we can see what's going on.
	InitLogging()

	// Make sure the XDG directories exist before we start writing to them.
	err := createXDGDirectories()
	if err != nil {
		log.Fatalf("config.init: %v", err)
	}

	// These values initialised here are required for the generation of the config file.
	InitP2P()
	InitHttp()

	if generateFlag() {

		// Reinit logging to STDOUT
		viper.Set("log.file", "STDOUT")
		InitLogging()

		switch Mode() {
		case "actor":
			log.Info("Generating new actor and node identity")
			actor, node := handleGenerateOrExit()
			generateActorConfigFile(actor, node)
		case "pong":
			viper.Set("actor.nick", "pong")
			log.Info("Generating new pong actor and node identity")
			actor, node := handleGenerateOrExit()
			generatePongConfigFile(actor, node)
		case "relay":
			viper.Set("actor.nick", "relay")
			log.Info("Generating new relay node identity")
			_, node := handleGenerateOrExit()
			generateRelayConfigFile(node)
		}
		os.Exit(0)
	}

	// Look in the current directory, the home directory and /etc for the config file.
	// In that order.
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(configHome)

	// We *must* read the config file after we have generated the identity.
	// Otherwise: Unforeseen consequences.
	log.Infof("Using config file: %s", configFile())
	viper.SetConfigFile(configFile())
	err = viper.ReadInConfig()
	if err != nil {
		log.Warnf("No config file found: %s", err)
	}

	if !RelayMode() {
		InitActor()
	}

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
		log.Fatalf("error: %v", err)
	}

	fmt.Println("# " + ActorKeyset().DID.Id)

	return fmt.Print(string(configYAML))
}

func Save() error {

	return viper.SafeWriteConfig()

}

func DataHome() string {
	return dataHome
}

func ConfigHome() string {
	return configHome
}

// Return the configName to use. If the mode is not the default, return the mode.
// If the mode is the default, return the actor nick provided on the command line.
// This is to avoid confusion for the user, so they don't have to use a nick
// to know which config file to use.
func configName() string {

	if Mode() != defaultMode {
		return Mode()
	}

	return pflag.CommandLine.Lookup("nick").Value.String()

}

// Returns the configfile name to use.
// The preferred value is the explcitily requested config file on the command line.
// Else it uses the nick of the actor or the mode.
func configFile() string {

	var (
		filename string
		err      error
	)

	// Prefer explcitly requested config. If not, use the name of the actor or mode.
	if config != defaultConfigFile {
		filename, err = homedir.Expand(config)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
	} else {
		filename = configHome + configName() + ".yaml"
	}

	return filepath.Clean(filename)

}

func createXDGDirectories() error {

	err := os.MkdirAll(configHome, configDirMode)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dataHome, dataHomeMode)
	if err != nil {
		return err
	}

	return nil

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

func showDefaultsFlag() bool {
	showDefaultsFlag, err := pflag.CommandLine.GetBool("show-defaults")
	if err != nil {
		log.Warnf("config.init: %v", err)
		return false
	}

	return showDefaultsFlag
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
