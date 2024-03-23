package config

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

const (
	NAME       string = "go-ma-actor"
	VERSION    string = "v0.3.1"
	ENV_PREFIX string = "GO_MA_ACTOR"

	configDirMode  os.FileMode = 0700
	configFileMode os.FileMode = 0600
	dataHomeMode   os.FileMode = 0755
)

type Config interface {
	MarshalToYAML() ([]byte, error)
	Print()
	Save() error
}

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
	if !GenerateFlag() {
		log.Infof("Using config file: %s", File()) // This one goes to STDERR
		viper.SetConfigFile(File())
		err = viper.ReadInConfig()
		if err != nil {
			log.Warnf("No config file found: %s", err)
		}
	}

	// Handle the easy flags first.
	if versionFlag() {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	// Make sure the XDG directories exist before we start writing to them.
	err = createXDGDirectories()
	if err != nil {
		log.Fatal(err)
	}

	return nil

}

func PrintAll() (int, error) {

	configMap := viper.AllSettings()

	configYAML, err := yaml.Marshal(configMap)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("# " + ActorKeyset().DID.Id)

	return fmt.Println(string(configYAML))
}

func Print(c Config) {
	configYAML, err := c.MarshalToYAML()
	if err != nil {
		log.Fatalf("Failed to marshal config to YAML: %v", err)
	}

	fmt.Println(string(configYAML))
}

// Write the generated config to the correct file
// NB! This fails fatally in case of an error.
func Save(c Config) error {
	filePath := File()
	var errMsg string

	// Determine the file open flags based on the forceFlag
	var flags int
	if ForceFlag() {
		// Allow overwrite
		log.Warnf("Force flag set, overwriting existing config file %s", filePath)
		flags = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	} else {
		// Prevent overwrite
		flags = os.O_WRONLY | os.O_CREATE | os.O_EXCL
	}

	file, err := os.OpenFile(filePath, flags, configFileMode)
	if err != nil {
		if os.IsExist(err) {
			errMsg = fmt.Sprintf("File %s already exists.", filePath)
		} else {
			errMsg = fmt.Sprintf("Failed to open file: %v", err)
		}
		return fmt.Errorf(errMsg)
	}
	defer file.Close()

	content, err := c.MarshalToYAML()
	if err != nil {
		return fmt.Errorf("failed to marshal to YAML: %w", err)
	}

	// Write content to file.
	if _, err := file.Write(content); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	log.Printf("Generated config file %s", filePath)
	return nil
}
