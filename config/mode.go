package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	DEFAULT_PONG_REPLY        = "Pong!"
	DEFAULT_PONG_FORTUNE_MODE = false
	DEFAULT_PONG_FORTUNE_ARGS = "-s"

	pong  = "pong"
	relay = "relay"

	defaultMode = "actor"
)

var ErrConflictingModes = "Can't have both pong and relay mode enabled at the same time."

// NB! This file is used early in the initialization process, so it can't depend on other packages.

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

		// This only sets fallback when not specified on command line
		SetProfile(pong)
		return pong
	}

	if RelayMode() {

		// This only sets fallback when not specified on command line
		SetProfile(relay)
		return relay
	}

	return defaultMode
}

func PongFortuneMode() bool {
	return viper.GetBool("mode.pong.fortune.enable") && PongMode()
}

func PongFortuneArgs() []string {
	args := viper.GetString("mode.pong.fortune.args")

	return []string{args}
}
