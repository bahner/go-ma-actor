package config

import "os"

// This should be done after flag.Parse() in main.
func Init() {
	InitLogging()
	InitNodeIdentity()
	InitIdentity()

	if *genenv {
		PrintEnvironment()
		os.Exit(0)
	}

}
