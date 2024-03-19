package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bahner/go-ma-actor/config"
	log "github.com/sirupsen/logrus"
)

const relay = "relay"

// Run the pong actor. Cancel it from outside to stop it.
func main() {

	ctx := context.Background()
	initConfig(relay)

	p, err := initP2P()
	if err != nil {
		fmt.Printf("Failed to initialize p2p: %v\n", err)
		return
	}

	go p.StartDiscoveryLoop(ctx)
	fmt.Println("Discovery loop started.")

	// Start a simple web server to handle incoming requests.
	// This is defined in web.go. It makes it possible to add extra parameters to the handler.
	mux := http.NewServeMux()
	h := &WebEntity{
		P2P: p,
	}
	mux.HandleFunc("/", h.WebHandler)

	log.Infof("Listening on %s", config.HttpSocket())

	// IN relay mode we want to stop here.
	fmt.Println("Web server starting on http://" + config.HttpSocket() + "/")
	http.ListenAndServe(config.HttpSocket(), mux)

}
