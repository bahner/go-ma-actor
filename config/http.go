package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	defaultHttpSocket  string = "127.0.0.1:5002"
	defaultHttpRefresh int    = 10
)

func init() {

	pflag.String("http-socket", defaultHttpSocket, "Address for webserver to listen on")
	pflag.Int("http-refresh", defaultHttpRefresh, "Number of seconds for webpages to wait before refresh")

	viper.BindPFlag("http.socket", pflag.Lookup("http-socket"))
	viper.BindPFlag("http.refresh", pflag.Lookup("http-refresh"))

	viper.SetDefault("http.socket", defaultHttpSocket)
	viper.SetDefault("http.refresh", defaultHttpRefresh)

}

func HttpSocket() string {
	return viper.GetString("http.socket")
}

func HttpRefresh() int {
	return viper.GetInt("http.refresh")
}
