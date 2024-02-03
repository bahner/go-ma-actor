package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	defaultHttpSocket string = "0.0.0.0:5003"
)

func init() {

	// Flags - user configurations

	pflag.String("http_socket", defaultHttpSocket, "Address for webserver to listen on")
	viper.BindPFlag("http.socket", pflag.Lookup("socket"))

}

func GetHttpSocket() string {
	return viper.GetString("http.socket")
}
