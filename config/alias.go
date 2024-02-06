package config

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const defaultAliases = "~/.ma/aliases.db"

func init() {

	pflag.String("aliases", defaultAliases, "File to *write* node aliases to. If the file does not exist, it will be created.")
	viper.BindPFlag("aliases", pflag.Lookup("aliases"))
	viper.SetDefault("aliases", defaultAliases)
}

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
