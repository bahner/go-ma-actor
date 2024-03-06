package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

const (
	defaultLogLevel string      = "info"
	logFilePerm     os.FileMode = 0640
)

var defaultLogfile string = NormalisePath(dataHome + defaultNick + ".log")

func init() {

	pflag.String("log-level", defaultLogLevel, "Loglevel to use for application.")
	pflag.String("log-file", defaultLogfile, "Logfile to use for application. Accepts 'STDERR' and 'STDOUT' as such.")
}

func InitLogging() {

	viper.BindPFlag("log.file", pflag.Lookup("log-file"))
	viper.SetDefault("log.file", genDefaultLogfile(profile))

	viper.BindPFlag("log.level", pflag.Lookup("log-level"))
	viper.SetDefault("log.level", defaultLogLevel)

	// Init logger
	ll, err := log.ParseLevel(viper.GetString("log.level"))
	if err != nil {
		fmt.Println(err)
		os.Exit(64) // EX_USAGE
	}
	log.SetLevel(ll)
	logfile, err := getLogFile()
	if err != nil {
		fmt.Println(err)
		os.Exit(73) // EX_CANTCREAT
	}
	if logfile == "STDERR" {
		log.SetOutput(os.Stderr)
	} else if logfile == "STDOUT" {
		log.SetOutput(os.Stdout)
	} else {
		file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, logFilePerm)
		if err != nil {
			fmt.Println(err)
			os.Exit(73) // EX_CANTCREAT
		}
		log.SetOutput(file)

	}
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})

	log := log.WithField("prefix", ActorNick())

	log.Info("Logger initialized with loglevel ", viper.GetString("log.level"), " and logfile ", viper.GetString("log.file"))

}

func getLogFile() (string, error) {
	lf := viper.GetString("log.file")

	if lf == "STDERR" || lf == "STDOUT" {
		return lf, nil
	}

	lf, err := homedir.Expand(lf)
	if err != nil {
		return "", err
	}

	lf, err = filepath.Abs(lf)
	if err != nil {
		return "", err
	}
	return NormalisePath(lf), nil
}

func LogLevel() string {
	return viper.GetString("log.level")
}

func LogFile() string {
	return viper.GetString("log.file")
}

func genDefaultLogfile(nick string) string {
	return NormalisePath(dataHome + nick + ".log")
}
