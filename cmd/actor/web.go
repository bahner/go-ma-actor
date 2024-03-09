package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/entity/actor"
	mactor "github.com/bahner/go-ma-actor/mode/actor"
	"github.com/bahner/go-ma-actor/mode/relay"

	"github.com/bahner/go-ma-actor/p2p"
	log "github.com/sirupsen/logrus"
)

// NB! In relay mode we stop here.
func startWebServer(p *p2p.P2P, a *actor.Actor) {

	fmt.Println("Starting web server...")

	// When this function stops the app stops.
	defer os.Exit(1)

	// Start a simple web server to handle incoming requests.
	// This is defined in web.go. It makes it possible to add extra parameters to the handler.
	mux := http.NewServeMux()

	// Different handlers for diiferent modes.
	if config.RelayMode() {
		h := &relay.WebHandlerData{
			P2P: p,
		}
		mux.HandleFunc("/", h.WebHandler)

	} else {
		h := &mactor.WebHandlerData{
			P2P:    p,
			Entity: a.Entity,
		}
		mux.HandleFunc("/", h.WebHandler)
	}

	log.Infof("Listening on %s", config.HttpSocket())

	// IN relay mode we want to stop here.
	fmt.Println("Web server starting on http://" + config.HttpSocket() + "/")
	err := http.ListenAndServe(config.HttpSocket(), mux)
	if err != nil {
		fmt.Println("failed.")
		// The webserver isn't critical, so we just log the error and continue.
		log.Errorf("Failed to start web server: %v", err)
		return
	}

}
