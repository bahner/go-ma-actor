package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	defaultHttpSocket string = "127.0.0.1:5002"
)

func init() {

	// Flags - user configurations

	pflag.String("http-socket", defaultHttpSocket, "Address for webserver to listen on")
	viper.BindPFlag("http.socket", pflag.Lookup("http-socket"))

}

func GetHttpSocket() string {
	return viper.GetString("http.socket")
}
