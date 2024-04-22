package config

var defaultProfile = "actor"

// Profile is the mode unless overridden by the profile flag.
func Profile() string {

	// This is used early so command line takes precedence
	if commonFlagset.Lookup("profile").Changed {
		return commonFlagset.Lookup("profile").Value.String()
	}

	return defaultProfile
}

// Call this for special cmds like "relay", "node", etc.
// Before initConfig() like calls.
func SetDefaultProfileName(p string) {
	defaultProfile = p
}
