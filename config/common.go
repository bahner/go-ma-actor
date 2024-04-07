package config

import (
	"fmt"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

var (
	CommonFlags = pflag.NewFlagSet("common", pflag.ContinueOnError)

	commonOnce sync.Once
)

func InitCommon() {

	commonOnce.Do(func() {

		// Allow to set config file via command line flag.
		CommonFlags.StringP("config", "c", "", "Config file to use.")
		CommonFlags.StringP("profile", "p", "", "Config profile (name) to use.")

		CommonFlags.Bool("show-config", false, "Whether to print the config.")

		CommonFlags.BoolP("version", "v", false, "Print version and exit.")

		CommonFlags.Bool("generate", false, "Generates a new keyset")
		CommonFlags.Bool("force", false, "Forces regneration of config keyset and publishing")

		CommonFlags.String("debug-socket", defaultDebugSocket, "Port to listen on for debug endpoints")

		if HelpNeeded() {
			fmt.Println("Common flags:")
			CommonFlags.PrintDefaults()
		} else {
			CommonFlags.Parse(os.Args[1:])
		}
	})
}

func GenerateFlag() bool {
	// This will exit when done. It will also publish if applicable.
	generateFlag, err := CommonFlags.GetBool("generate")
	if err != nil {
		log.Warnf("config.init: %v", err)
		return false
	}

	return generateFlag
}

func PublishFlag() bool {
	publishFlag, err := CommonFlags.GetBool("publish")
	if err != nil {
		log.Warnf("config.init: %v", err)
		return false
	}

	return publishFlag
}

func ShowConfigFlag() bool {
	showConfigFlag, err := CommonFlags.GetBool("show-config")
	if err != nil {
		log.Warnf("config.init: %v", err)
		return false
	}

	return showConfigFlag
}

func versionFlag() bool {
	versionFlag, err := CommonFlags.GetBool("version")
	if err != nil {
		log.Warnf("config.init: %v", err)
		return false
	}

	return versionFlag
}

func ForceFlag() bool {
	forceFlag, err := CommonFlags.GetBool("force")
	if err != nil {
		log.Warnf("config.init: %v", err)
		return false
	}

	return forceFlag
}

/*
	Parse common flags.

This is idemPotent in the sense that it can be called
multiple times without side effects, as the flags are
only parsed once.
Set exitOnHelp to true if you want the program to exit
after help is printed. This is useful for the main function,
when this is the last flag parsing function called.
*/
func ParseCommonFlags(exitOnHelp bool) {

	InitCommon()
	InitLog()
	InitDB()
	InitP2P()
	InitHTTP()

	if HelpNeeded() && exitOnHelp {
		os.Exit(0)
	}
}
