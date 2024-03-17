package config

import "github.com/spf13/pflag"

// NB! This file is used early in the initialization process, so it can't depend on other packages.
const defaultProfile string = "actor"

var profile string = defaultProfile

// Profile is the mode unless overridden by the profile flag.
func Profile() string {

	flag := pflag.Lookup("profile")
	if flag != nil && flag.Changed {
		return flag.Value.String()
	}

	return profile
}

func SetProfile(p string) {
	profile = p
}
