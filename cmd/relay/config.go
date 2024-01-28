package main

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// defaultListenPort int    = 4001 // 0 = random
	defaultHttpSocket string = "0.0.0.0:4000"
)

func init() {

	// Flags - user configurations

	pflag.String("http_socket", defaultHttpSocket, "Address for webserver to listen on")
	viper.SetDefault("http.socket", defaultHttpSocket)
	viper.BindPFlag("http.socket", pflag.Lookup("socket"))

	pflag.Parse()
}

func getHttpSocket() string {
	return viper.GetString("http.socket")
}
