package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	DefaultPongReply = "Pong!"

	defaultPongMode  = false
	defaultRelayMode = false
	defaultDebugMode = false

	defaultMode = "actor"
)

func init() {
	pflag.Bool("pong", defaultPongMode, "Pong mode with automatic replies and no UI.")
	viper.BindPFlag("mode.pong.enabled", pflag.Lookup("pong"))

	pflag.Bool("relay", defaultRelayMode, "The relay to use for the p2p network.")
	viper.BindPFlag("mode.relay", pflag.Lookup("relay"))

	pflag.Bool("debug", defaultDebugMode, "Enable debug mode. Adds pprof and other debug endpoints to the webserver.")
	viper.BindPFlag("mode.debug", pflag.Lookup("debug"))

}

// Returns the mode of the actor as as a string, eg. "actor", "pong", "relay".
func InitMode() string {

	if PongMode() && RelayMode() {
		log.Fatal("Can't have both pong and relay mode enabled at the same time.")
	}

	if log.GetLevel() == log.DebugLevel {
		log.Info("Debug mode enabled due to loglevel.")
		viper.Set("mode.debug", true)
	}

	if PongMode() {
		return "pong"
	}

	if RelayMode() {
		return "relay"
	}

	return defaultMode
}

func DebugMode() bool {
	return viper.GetBool("mode.debug")
}

// If actor.home is set to pong, then we are in pong mode.
// THIs means that we don't render the ui and reply automatically to messages.
func PongMode() bool {

	// if GetHome() == pongTriggerHomeName {
	// 	return true
	// }

	return viper.GetBool("mode.pong.enabled")
}

func PongReply() string {
	return viper.GetString("mode.pong.reply")
}

func RelayMode() bool {
	return viper.GetBool("mode.relay")
}
