package config

import (
	"fmt"
	"sync"

	"github.com/spf13/pflag"
)

var (
	commonFlagset = pflag.NewFlagSet("common", pflag.ExitOnError)

	commonFlagsOnce sync.Once

	debugFlag bool = false
	forceFlag bool = false

	generateCommandFlag   bool = false
	showConfigCommandFlag bool = false
	versionCommandFlag    bool = false
)

func InitCommonFlagset() {

	commonFlagsOnce.Do(func() {

		// Allow to set config file via command line flag.
		commonFlagset.StringP("config", "c", "", "Config file to use.")
		commonFlagset.StringP("profile", "p", "", "Config profile (name) to use.")

		// Commands
		commonFlagset.BoolVar(&showConfigCommandFlag, "show-config", false, "Whether to print the config.")
		commonFlagset.BoolVarP(&versionCommandFlag, "version", "v", false, "Print version and exit.")
		commonFlagset.BoolVar(&generateCommandFlag, "generate", false, "Generates a new keyset")

		// Flags
		commonFlagset.BoolVar(&forceFlag, "force", false, "Forces regneration of config keyset and publishing")
		commonFlagset.BoolVar(&debugFlag, "debug", false, "Port to listen on for debug endpoints")

		if HelpNeeded() {
			fmt.Println("Common Flags:")
			commonFlagset.PrintDefaults()
		}

	})
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

func Debug() bool {
	return debugFlag
}

func ForceFlag() bool {
	return forceFlag
}

func GenerateFlag() bool {
	return generateCommandFlag
}

func ShowConfigFlag() bool {
	return showConfigCommandFlag
}

func VersionFlag() bool {
	return versionCommandFlag
}
