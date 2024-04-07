package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/spf13/pflag"
)

var (
	CommonFlags = pflag.NewFlagSet("common", pflag.ContinueOnError)

	commonOnce sync.Once

	debugFlag bool = false
	forceFlag bool = false

	generateCommandFlag   bool = false
	showConfigCommandFlag bool = false
	versionCommandFlag    bool = false
)

func InitCommon() {

	commonOnce.Do(func() {

		// Allow to set config file via command line flag.
		CommonFlags.StringP("config", "c", "", "Config file to use.")
		CommonFlags.StringP("profile", "p", "", "Config profile (name) to use.")

		// COmmands
		CommonFlags.BoolVar(&showConfigCommandFlag, "show-config", false, "Whether to print the config.")
		CommonFlags.BoolVarP(&versionCommandFlag, "version", "v", false, "Print version and exit.")
		CommonFlags.BoolVar(&generateCommandFlag, "generate", false, "Generates a new keyset")

		// Flags
		CommonFlags.BoolVar(&forceFlag, "force", false, "Forces regneration of config keyset and publishing")
		CommonFlags.BoolVar(&debugFlag, "debug", false, "Port to listen on for debug endpoints")

		if HelpNeeded() {
			fmt.Println("Common flags:")
			CommonFlags.PrintDefaults()
		} else {
			CommonFlags.Parse(os.Args[1:])
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
