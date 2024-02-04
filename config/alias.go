package config

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// Returns expanded path to the aliases file
// If the expansion fails it returns an empty string
func GetAliases() string {

	path := viper.GetString("aliases")
	path, err := homedir.Expand(path)
	if err != nil {
		return ""
	}

	return path

}
