package main

import (
	"github.com/bahner/go-ma-actor/config"
	"github.com/spf13/pflag"
)

const (
	defaultPongMode  = false
	defaultRelayMode = false
)

func init() {

	// NB! Other mode pflags are in the proper mode packages.
	pflag.Bool("pong", defaultPongMode, "Pong mode with automatic replies and no UI.")
	pflag.Bool("pong-fortune", config.DEFAULT_PONG_FORTUNE_MODE, "Reply with a fortune cookie, instead of a static message, if availble.")
	pflag.Bool("relay", defaultRelayMode, "Relay mode with no actor, to just listen and relay messages.")

}
