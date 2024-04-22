package config

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

func ActorFlagset() *pflag.FlagSet {

	mergeActorFlagset()

	return actorFlagset
}

func ActorFlagsetParse(exitOnHelp bool) {

	mergeActorFlagset()

	if HelpNeeded() && exitOnHelp {
		os.Exit(0)
	}

	err = actorFlagset.Parse(os.Args[1:])
	if err != nil {
		fmt.Printf("Error parsing actor flags: %s", err)
		os.Exit(64) // EX_USAGE
	}
}

func CommonFlagset() *pflag.FlagSet {

	mergeCommonFlagset()

	return commonFlagset
}

func CommonFlagsetParse(exitOnHelp bool) {

	mergeCommonFlagset()
	mergeFromFlagsetInto(commonFlagset, actorFlagset)

	if HelpNeeded() && exitOnHelp {
		os.Exit(0)
	}

	err = commonFlagset.Parse(os.Args[1:])
	if err != nil {
		fmt.Printf("Error parsing common flags: %s", err)
		os.Exit(64) // EX_USAGE
	}
}

func mergeActorFlagset() {
	mergeCommonFlagset()
	actorFlags()

	mergeFromFlagsetInto(commonFlagset, actorFlagset)

}

func mergeCommonFlagset() {

	InitCommonFlagset()
	initLogFlagset()
	initDBFlagset()
	initP2PFlagset()
	initHTTPFlagset()

	mergeFromFlagsetInto(logFlagset, commonFlagset)
	mergeFromFlagsetInto(dbFlagset, commonFlagset)
	mergeFromFlagsetInto(p2pFlagset, commonFlagset)
	mergeFromFlagsetInto(httpFlagset, commonFlagset)

}

// Function to add flags from one set to another
func mergeFromFlagsetInto(from, to *pflag.FlagSet) {
	from.VisitAll(func(flag *pflag.Flag) {
		to.AddFlag(flag)
	})
}
