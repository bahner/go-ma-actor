package main

import (
	"flag"

	"go.deanishe.net/env"
)

const (
	defaultListenPort string = "4001" // 0 = random
	defaultHttpAddr   string = "0.0.0.0"
	defaultHttpPort   string = "4000"
)

var (
	httpSocket string

	httpAddr   string = env.Get("GO_MA_RELAY_HTTP_ADDR", defaultHttpAddr)
	httpPort   string = env.Get("GO_MA_RELAY_HTTP_PORT", defaultHttpPort)
	listenPort string = env.Get("GO_MA_RELAY_LISTEN_PORT", defaultListenPort)
)

func init() {

	// Flags - user configurations

	flag.StringVar(&httpAddr, "httpAddr", httpAddr, "Address to listen on")
	flag.StringVar(&httpPort, "httpPort", httpPort, "Listen port for webserver")

	flag.StringVar(&listenPort, "listenPort", listenPort, "Port to listen on for peers")

	flag.Parse()

	// Assemble vars for http server
	httpSocket = httpAddr + ":" + httpPort
}

func GetListenPort() string {
	return listenPort
}

func GetHttpSocket() string {
	return httpSocket
}

func GetHttpAddr() string {
	return httpAddr
}

func GetHttpPort() string {
	return httpPort
}
