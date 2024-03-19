package web

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bahner/go-ma-actor/config"

	log "github.com/sirupsen/logrus"
)

// Start the WebU with the given handler.
func Start(h WebHandler) {

	fmt.Println("Starting web server...")

	// When this function stops the app stops.
	defer os.Exit(1)

	// Start a simple web server to handle incoming requests.
	// This is defined in web.go. It makes it possible to add extra parameters to the handler.
	mux := http.NewServeMux()
	mux.Handle("/", h)

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
