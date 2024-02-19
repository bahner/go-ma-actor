package config

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

const (
	defaultLogLevel string = "info"
	defaultLogfile  string = "." + NAME + ".log"
)

func init() {

	pflag.String("log-level", defaultLogLevel, "Loglevel to use for application.")
	viper.SetDefault("log.level", defaultLogLevel)
	viper.BindPFlag("log.level", pflag.Lookup("log-level"))

	pflag.String("log-file", defaultLogfile, "Logfile to use for application. Accepts 'STDERR' and 'STDOUT' as such.")
	viper.SetDefault("log.file", defaultLogfile)
	viper.BindPFlag("log.file", pflag.Lookup("log-file"))
}

func InitLogging() {

	// Init logger
	ll, err := log.ParseLevel(viper.GetString("log.level"))
	if err != nil {
		fmt.Println(err)
		os.Exit(64) // EX_USAGE
	}
	log.SetLevel(ll)
	logfile := viper.GetString("log.file")
	if logfile == "STDERR" {
		log.SetOutput(os.Stderr)
	} else if logfile == "STDOUT" {
		log.SetOutput(os.Stdout)
	} else {
		logfile, err = homedir.Expand(logfile)
		if err != nil {
			fmt.Println(err)
			os.Exit(73) // EX_CANTCREAT
		}
		file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			fmt.Println(err)
			os.Exit(73) // EX_CANTCREAT
		}
		log.SetOutput(file)

	}

	log.Info("Logger initialized with loglevel ", viper.GetString("log.level"), " and logfile ", viper.GetString("log.file"))

}

func GetLogLevel() string {
	return viper.GetString("log.level")
}

func GetLogFile() string {
	return viper.GetString("log.file")
}
