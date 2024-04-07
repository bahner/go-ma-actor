package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	defaultHttpSocket  string = "127.0.0.1:5002"
	defaultHttpRefresh int    = 10
)

var (
	httpFlags = pflag.NewFlagSet("http", pflag.ContinueOnError)
	httpOnce  sync.Once
)

type HTTPConfig struct {
	Socket      string `yaml:"socket"`
	Refresh     int    `yaml:"refresh"`
	DebugSocket string `yaml:"debug_socket"`
}

func InitHTTP() {

	httpOnce.Do(func() {
		httpFlags.String("http-socket", defaultHttpSocket, "Address for webserver to listen on")
		httpFlags.Int("http-refresh", defaultHttpRefresh, "Number of seconds for webpages to wait before refresh")

		viper.BindPFlag("http.socket", httpFlags.Lookup("http-socket"))
		viper.BindPFlag("http.refresh", httpFlags.Lookup("http-refresh"))

		viper.SetDefault("http.socket", defaultHttpSocket)
		viper.SetDefault("http.refresh", defaultHttpRefresh)

		if HelpNeeded() {
			fmt.Println("HTTP Flags:")
			httpFlags.PrintDefaults()
		} else {
			httpFlags.Parse(os.Args[1:])
		}
	})
}

func HTTP() HTTPConfig {

	return HTTPConfig{
		Socket:      HttpSocket(),
		Refresh:     HttpRefresh(),
		DebugSocket: HttpDebugSocket(),
	}
}

func HttpSocket() string {
	return viper.GetString("http.socket")
}

func HttpRefresh() int {
	return viper.GetInt("http.refresh")
}

func HttpDebugSocket() string {
	return viper.GetString("http.debug-socket")
}
