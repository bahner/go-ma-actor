package main

import (
	"errors"
	"os"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma/did/doc"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	defaultPongReply   = "Pong!"
	defaultFortuneMode = false
	pong               = "pong"
)

var defaultFortuneArgs = []string{"-s"}

func init() {
	pflag.String("pong-reply", defaultPongReply, "The message to send back to the sender")

	viper.BindPFlag("mode.pong.reply", pflag.Lookup("pong-reply"))
	viper.SetDefault("mode.pong.reply", defaultPongReply)

	viper.BindPFlag("mode.pong.fortune.enable", pflag.Lookup("pong-fortune"))
	viper.SetDefault("mode.pong.fortune.enable", defaultFortuneMode)

	viper.SetDefault("mode.pong.fortune.args", defaultFortuneArgs)
}

func pongFortuneMode() bool {
	return viper.GetBool("mode.pong.fortune.enable")
}

func pongFortuneArgs() []string {
	return viper.GetStringSlice("mode.pong.fortune.args")
}

func pongReply() string {
	return viper.GetString("mode.pong.reply")
}

func initConfig(profile string) {

	// Always parse the flags first
	config.InitCommonFlags()
	config.InitActorFlags()
	pflag.Parse()
	config.SetProfile(profile)
	config.Init()

	if config.GenerateFlag() {
		// Reinit logging to STDOUT
		log.SetOutput(os.Stdout)
		log.Info("Generating new actor and node identity")
		actor, node := generateActorIdentitiesOrPanic(pong)
		actorConfig := configTemplate(actor, node)
		config.Generate(actorConfig)
		os.Exit(0)
	}

	config.InitActor()

	// This flag is dependent on the actor to be initialized to make sense.
	if config.ShowConfigFlag() {
		config.Print()
		os.Exit(0)
	}

}

func generateActorIdentitiesOrPanic(name string) (string, string) {
	actor, node, err := config.GenerateActorIdentities(name)
	if err != nil {
		if errors.Is(err, doc.ErrAlreadyPublished) {
			log.Warnf("Actor document already published: %v", err)
		} else {
			log.Fatal(err)
		}
	}
	return actor, node
}
