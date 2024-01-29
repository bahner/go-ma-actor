package main

import (
	"context"
	"net/http"
	"os"

	"github.com/bahner/go-ma-actor/actor"
	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/p2p"
	"github.com/bahner/go-ma-actor/ui"
	"github.com/spf13/pflag"

	log "github.com/sirupsen/logrus"
)

func main() {

	pflag.Parse()

	config.Init("actor")
	config.InitLogging()
	config.InitP2P()
	config.InitActor()

	p, err := p2p.Init(nil)
	if err != nil {
		log.Errorf("failed to initialize p2p: %v", err)
		os.Exit(75)
	}

	// We need to discover peers before we can do anything else.
	// So this is a blocking call.
	p.DiscoverPeers()

	if err != nil {
		log.Errorf("failed to initialize p2p: %v", err)
		os.Exit(75)
	}

	// Now we can start continuous discovery in the background.
	go p.DiscoveryLoop(context.Background())

	a, err := actor.NewFromKeyset(config.GetKeyset(), config.GetPublish())
	if err != nil || a == nil {
		log.Errorf("failed to create actor: %v", err)
		os.Exit(70)
	}

	// Start a simple web server to handle incoming requests.
	// This is defined in web.go. It makes it possible to add extra parameters to the handler.
	h := &WebHandlerData{p.Node, a}
	http.HandleFunc("/", h.WebHandler)
	log.Infof("Listening on %s\n", config.GetHttpSocket())
	go http.ListenAndServe(config.GetHttpSocket(), nil)

	e := config.GetHome()
	// Draw the UI.
	log.Debugf("Starting text UI")
	ui := ui.NewChatUI(p, a, e)
	if err := ui.Run(); err != nil {
		log.Errorf("error running text UI: %s", err)
	}
}
