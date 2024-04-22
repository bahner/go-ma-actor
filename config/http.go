package config

import (
	"fmt"
	"sync"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	defaultHttpSocket  string = "127.0.0.1:5002"
	defaultHttpRefresh int    = 10
)

var (
	httpFlagset   = pflag.NewFlagSet("http", pflag.ExitOnError)
	httpFlagsOnce sync.Once
)

type HTTPConfig struct {
	Socket      string `yaml:"socket"`
	Refresh     int    `yaml:"refresh"`
	DebugSocket string `yaml:"debug_socket"`
}

func initHTTPFlagset() {

	httpFlagsOnce.Do(func() {
		httpFlagset.String("http-socket", defaultHttpSocket, "Address for webserver to listen on")
		httpFlagset.Int("http-refresh", defaultHttpRefresh, "Number of seconds for webpages to wait before refresh")

		viper.BindPFlag("http.socket", httpFlagset.Lookup("http-socket"))
		viper.BindPFlag("http.refresh", httpFlagset.Lookup("http-refresh"))

		viper.SetDefault("http.socket", defaultHttpSocket)
		viper.SetDefault("http.refresh", defaultHttpRefresh)

		if HelpNeeded() {
			fmt.Println("HTTP Flags:")
			httpFlagset.PrintDefaults()
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
