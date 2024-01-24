package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

const (
	defaultLogLevel string = "info"
	defaultLogfile  string = NAME + ".log"
)

func InitLogging() {

	// Init logger
	ll, err := log.ParseLevel(viper.GetString("log.level"))
	if err != nil {
		fmt.Println(err)
		os.Exit(64) // EX_USAGE
	}
	log.SetLevel(ll)
	file, err := os.OpenFile(viper.GetString("log.file"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(73) // EX_CANTCREAT
	}
	log.SetOutput(file)

	log.Info("Logger initialized with loglevel ", viper.GetString("log.level"), " and logfile ", viper.GetString("log.file"))

}

func GetLogLevel() string {
	return viper.GetString("log.level")
}

func GetLogFile() string {
	return viper.GetString("log.file")
}
