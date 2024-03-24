package config

import "github.com/spf13/pflag"

var defaultProfile = "actor"

// Profile is the mode unless overridden by the profile flag.
func Profile() string {

	flag := pflag.Lookup("profile")
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
