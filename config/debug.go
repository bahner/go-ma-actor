//go:build debug

package config

import (
	"net/http"
	_ "net/http/pprof"
	"runtime"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const defaultDebugSocket = "0.0.0.0:6060"

func init() {
	// Assume you have a function setupDebugHandlers to register debug routes
	setupDebugHandlers()

	pflag.String("debug-socket", defaultDebugSocket, "Port to listen on for debug endpoints")
	viper.BindPFlag("http.debug-socket", pflag.Lookup("debug-socket"))
	viper.SetDefault("http.debug-socket", defaultDebugSocket)
}

func setupDebugHandlers() {
	// Register your pprof handlers or other debug routes here
	// Since "net/http/pprof" is imported above, its init function automatically registers its routes with the default mux
	go http.ListenAndServe(HttpDebugSocket(), nil)

	http.HandleFunc("/force-gc", func(w http.ResponseWriter, r *http.Request) {
		// Force a garbage collection
		runtime.GC()
		w.Write([]byte("Garbage collection triggered"))
	})
}
