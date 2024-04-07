package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/bahner/go-ma-actor/internal"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

const (
	defaultLogLevel string      = "info"
	logFilePerm     os.FileMode = 0640
)

var (
	logFlags = pflag.NewFlagSet("log", pflag.ContinueOnError)
	logOnce  sync.Once
)

func init() {
	os.Setenv("GOLOG_OUTPUT", "")
}

func InitLog() {

	logOnce.Do(func() {
		logFlags.String("loglevel", defaultLogLevel, "Loglevel to use for application.")
		logFlags.String("logfile", defaultLogfile(), "Logfile to use for application. Accepts 'stderr' and 'stdout' as such.")

		viper.BindPFlag("log.file", logFlags.Lookup("logfile"))
		viper.BindPFlag("log.level", logFlags.Lookup("loglevel"))

		viper.SetDefault("log.level", defaultLogLevel)
		viper.SetDefault("log.file", defaultLogfile)

		if HelpNeeded() {
			fmt.Println("Log Flags:")
			logFlags.PrintDefaults()
		} else {
			logFlags.Parse(os.Args[1:])
		}

	})
}

type LogConfig struct {
	Level string `yaml:"level"`
	File  string `yaml:"file"`
}

func Log() LogConfig {

	viper.SetDefault("log.file", genDefaultLogFileName(Profile()))

	initLogging()

	return LogConfig{
		Level: LogLevel(),
		File:  LogFile(),
	}
}

func LogLevel() string {
	return viper.GetString("log.level")
}

func LogFile() string {
	return viper.GetString("log.file")
}

func defaultLogfile() string {
	return internal.NormalisePath(dataHome + Profile() + ".log")
}

func initLogging() {

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
	os.Setenv("GOLOG_FILE", logfile)
	if logfile == "stderr" {
		log.SetOutput(os.Stderr)
		os.Setenv("GOLOG_OUTPUT", logfile)
	} else if logfile == "stdout" {
		log.SetOutput(os.Stdout)
		os.Setenv("GOLOG_OUTPUT", logfile)
	} else {
		file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, logFilePerm)
		if err != nil {
			fmt.Println(err)
			os.Exit(73) // EX_CANTCREAT
		}
		log.SetOutput(file)
		os.Setenv("GOLOG_OUTPUT", "file")

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

	if lf == "stderr" || lf == "stdout" {
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
	return internal.NormalisePath(lf), nil
}

func genDefaultLogFileName(name string) string {
	return internal.NormalisePath(dataHome + name + ".log")
}
