package main

import (
	"os"

	"github.com/bahner/go-ma-actor/config"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	name = "robot"
)

func init() {
	pflag.String("openai-key", "", "The (paid) key to use with the OpenAI API")

	viper.BindPFlag("mode.openai.key", pflag.Lookup("openai-key"))
}

func initConfig(profile string) {

	// Always parse the flags first
	config.ActorFlags()
	pflag.Parse()
	config.SetProfile(profile)
	config.Init()


	config.InitActor()

	// This flag is dependent on the actor to be initialized to make sense.
	if config.ShowConfigFlag() {
		config.Print()
		os.Exit(0)
	}

}
func openAIKey() string {
	return viper.GetString("mode.openai.key")
}


