package config

import (
	"os"
)

var helpNeeded bool = false

func init() {
	for _, arg := range os.Args[1:] {
		if arg == "-h" || arg == "--help" {
			helpNeeded = true
			return
		}
	}
}

func HelpNeeded() bool {
	return helpNeeded
}
