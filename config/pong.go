package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	pongTriggerHomeName = "pong"
	defaultPongMode     = false
)

func init() {
	pflag.Bool("pong-mode", defaultPongMode, "Pong mode with automatic replies and no UI. Can also be triggere by setting actor.home to '"+pongTriggerHomeName+"'")
	viper.BindPFlag("pong.mode", pflag.Lookup("pong"))
}

// If actor.home is set to pong, then we are in pong mode.
// THIs means that we don't render the ui and reply automatically to messages.
func PongMode() bool {

	if GetHome() == pongTriggerHomeName {
		return true
	}

	return viper.GetBool("pong.mode")
}
