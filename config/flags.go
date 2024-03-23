package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

func init() {

	// Allow to set config file via command line flag.
	pflag.StringP("config", "c", "", "Config file to use.")
	pflag.StringP("profile", "p", "", "Config profile (name) to use.")

	pflag.Bool("show-config", false, "Whether to print the config.")
	pflag.BoolP("version", "v", false, "Print version and exit.")
	pflag.Bool("generate", false, "Generates a new keyset")
	pflag.Bool("publish", false, "Publishes keyset to IPFS")
	pflag.Bool("force", false, "Forces regneration of config keyset and publishing")
}

func GenerateFlag() bool {
	// This will exit when done. It will also publish if applicable.
	generateFlag, err := pflag.CommandLine.GetBool("generate")
	if err != nil {
		log.Warnf("config.init: %v", err)
		return false
	}

	return generateFlag
}

func PublishFlag() bool {
	publishFlag, err := pflag.CommandLine.GetBool("publish")
	if err != nil {
		log.Warnf("config.init: %v", err)
		return false
	}

	return publishFlag
}

func ShowConfigFlag() bool {
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

func ForceFlag() bool {
	forceFlag, err := pflag.CommandLine.GetBool("force")
	if err != nil {
		log.Warnf("config.init: %v", err)
		return false
	}

	return forceFlag
}
