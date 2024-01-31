package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/bahner/go-ma-actor/alias"
	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma-actor/p2p"
	"github.com/bahner/go-ma-actor/ui"
	"github.com/bahner/go-ma/did"
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

	a, err := entity.NewFromKeyset(config.GetKeyset(), config.GetKeyset().DID.Fragment)
	if err != nil || a == nil {
		log.Errorf("failed to create actor: %v", err)
		os.Exit(70)
	}

	eas := alias.GetEntityAliases()
	for _, ea := range eas {
		fmt.Printf("Entity alias: %s %s\n", ea.Nick, ea.Did)
	}
	// na := config.GetNodeAliases()

	// Start a simple web server to handle incoming requests.
	// This is defined in web.go. It makes it possible to add extra parameters to the handler.
	h := &WebHandlerData{p.Node, a}
	http.HandleFunc("/", h.WebHandler)
	log.Infof("Listening on %s", config.GetHttpSocket())
	go http.ListenAndServe(config.GetHttpSocket(), nil)

	home, err := did.New(config.GetHome())
	if err != nil {
		log.Errorf("home is not a valid DID: %v", err)
		os.Exit(70)
	}

	e, err := entity.New(home, nil, home.Fragment)
	if err != nil {
		log.Errorf("home is not a valid entity: %v", err)
		os.Exit(70)
	}
	// Draw the UI.
	log.Debugf("Starting text UI")
	ui, err := ui.NewChatUI(p, a, e)
	if err != nil {
		log.Errorf("error creating text UI: %s", err)
		os.Exit(75)
	}
	if err := ui.Run(); err != nil {
		log.Errorf("error running text UI: %s", err)
	}
}
