//go:build debug

package main

import (
	"net/http"
	_ "net/http/pprof"
	"runtime"
)

func init() {
	// Assume you have a function setupDebugHandlers to register debug routes
	setupDebugHandlers()
}

func setupDebugHandlers() {
	// Register your pprof handlers or other debug routes here
	// Since "net/http/pprof" is imported above, its init function automatically registers its routes with the default mux
	go http.ListenAndServe("localhost:6060", nil)

	http.HandleFunc("/force-gc", func(w http.ResponseWriter, r *http.Request) {
		// Force a garbage collection
		runtime.GC()
		w.Write([]byte("Garbage collection triggered"))
	})
}
