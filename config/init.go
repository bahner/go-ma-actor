package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// This should be done after flag.Parse() in main.
func Init() {

	if viper.GetBool("version") {
		fmt.Println(GetVersion())
		os.Exit(0)
	}

	InitLogging()
	InitNodeIdentity()
	InitIdentity()
}
