package config

import "github.com/spf13/pflag"

// NB! This file is used early in the initialization process, so it can't depend on other packages.

// Profile is the mode unless overridden by the profile flag.
func Profile() string {

	flag := pflag.Lookup("profile")
	if flag != nil {
		if flag.Changed {
			return flag.Value.String()
		}
	}

	return Mode()
}
