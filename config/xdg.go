package config

import (
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/bahner/go-ma"
	"github.com/bahner/go-ma-actor/internal"
	"github.com/mitchellh/go-homedir"
)

var (
	configHome string = xdg.ConfigHome + "/" + ma.NAME + "/"
	dataHome   string = xdg.DataHome + "/" + ma.NAME + "/"
)

// Returns the configfile name to use.
// The preferred value is the explcitily requested config file on the command line.
// Else it uses the nick of the actor or the mode.
func File() string {

	var (
		filename string
		err      error
	)

	config, err := CommonFlags.GetString("config")
	if err != nil {
		panic(err)
	}

	// Prefer explicitly requested config. If not, use the name of the profile name.
	if config != defaultConfigFile() && config != "" {
		filename, err = homedir.Expand(config)
		if err != nil {
			panic(err)
		}
	} else {
		filename = configHome + Profile() + ".yaml"
	}

	return filepath.Clean(filename)

}

func defaultConfigFile() string {
	return internal.NormalisePath(configHome + Profile() + ".yaml")
}

func XDGConfigHome() string {
	return configHome
}

func XDGDataHome() string {
	return dataHome
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
