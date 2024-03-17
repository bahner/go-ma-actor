package config

import (
	"github.com/spf13/pflag"
)

func InitCommonFlags() {

	// Allow to set config file via command line flag.
	pflag.StringP("config", "c", "", "Config file to use.")
	pflag.StringP("profile", "p", "", "Config profile (name) to use.")

	pflag.Bool("show-config", false, "Whether to print the config.")

	pflag.BoolP("version", "v", false, "Print version and exit.")

	pflag.Bool("generate", false, "Generates a new keyset")
	pflag.Bool("publish", false, "Publishes keyset to IPFS")
	pflag.Bool("force", false, "Forces regneration of config keyset and publishing")

	InitLogFlags()
	InitHTTPFlags()
	InitP2PFlags()
}
