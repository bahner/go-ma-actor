package config

import (
	"flag"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"go.deanishe.net/env"
)

var (

	// Logging
	logLevel = flag.String("loglevel", env.Get(GO_MA_ACTOR_LOGLEVEL_VAR, defaultLogLevel),
		"Loglevel to use for application. You can use environment variable "+GO_MA_ACTOR_LOGLEVEL_VAR+" to set this.")
	logfile = flag.String("logfile", env.Get(GO_MA_ACTOR_LOGFILE_VAR, defaultLogfile),
		"Logfile to use for application. You can use environment variable "+GO_MA_ACTOR_LOGFILE_VAR+" to set this.")
)

func InitLogging() {

	// Init logger
	level, err := log.ParseLevel(*logLevel)
	if err != nil {
		fmt.Println(err)
		os.Exit(64) // EX_USAGE
	}
	log.SetLevel(level)
	file, err := os.OpenFile(*logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(73) // EX_CANTCREAT
	}
	log.SetOutput(file)

	log.Info("Logger initialized")

}

func GetLogLevel() string {
	return *logLevel
}

func GetLogFile() string {
	return *logfile
}
