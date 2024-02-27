package config

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const defaultAliasesDB = "~/.ma/aliases.db"

func init() {

	pflag.String("db-aliases", defaultAliasesDB, "File to *write* node aliases to. If the file does not exist, it will be created.")
	viper.BindPFlag("db.aliases", pflag.Lookup("db-aliases"))
	viper.SetDefault("db.aliases", defaultAliasesDB)
}

// Returns expanded path to the aliases file
// If the expansion fails it returns an empty string
func GetAliasesDB() string {

	path := viper.GetString("aliases")
	path, err := homedir.Expand(path)
	if err != nil {
		return ""
	}

	return path

}
