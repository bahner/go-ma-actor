package config

import (
	log "github.com/sirupsen/logrus"
)

var defaultProfile = "actor"

// Profile is the mode unless overridden by the profile flag.
func Profile() string {

	flag := CommonFlags.Lookup("profile")
	log.Debugf("config.Profile: Lookup profile: %v", flag)
	if flag != nil && flag.Changed {
		return flag.Value.String()
	}

	return defaultProfile
}

// Call this for special cmds like "relay", "node", etc.
// Before initConfig() like calls.
func SetDefaultProfileName(p string) {
	defaultProfile = p
}
