package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	DEFAULT_PONG_REPLY = "Pong!"

	defaultPongMode  = false
	defaultRelayMode = false

	pong  = "pong"
	relay = "relay"

	defaultMode = "actor"
)

var ErrConflictingModes = "Can't have both pong and relay mode enabled at the same time."

// NB! This file is used early in the initialization process, so it can't depend on other packages.

func init() {
	// NB! Other mode pflags are in the proper mode packages.
	pflag.Bool("pong", defaultPongMode, "Pong mode with automatic replies and no UI.")
	pflag.Bool("relay", defaultRelayMode, "Relay mode with no actor, to just listen and relay messages.")
}

// If actor.home is set to pong, then we are in pong mode.
// THIs means that we don't render the ui and reply automatically to messages.
func PongMode() bool {

	pongMode, err := pflag.CommandLine.GetBool("pong")
	if err != nil {
		return false
	}

	return pongMode
}

func PongReply() string {
	return viper.GetString("mode.pong.reply")
}

func RelayMode() bool {

	relayMode, err := pflag.CommandLine.GetBool("relay")
	if err != nil {
		return false
	}

	return relayMode
}

// Returns the mode of the actor as as a string, eg. "actor", "pong", "relay".
func Mode() string {

	if PongMode() && RelayMode() {
		panic(ErrConflictingModes)
	}

	if PongMode() {
		return pong
	}

	if RelayMode() {
		return relay
	}

	return defaultMode
}
