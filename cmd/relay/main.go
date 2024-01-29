package main

import (
	"context"
	"net/http"

	"github.com/spf13/pflag"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/p2p"

	libp2p "github.com/libp2p/go-libp2p"
	log "github.com/sirupsen/logrus"
)

var (
	p *p2p.P2P
)

func main() {

	pflag.Parse()

	config.Init("relay")
	config.InitLogging()
	config.InitP2P()

	var err error

	ctx := context.Background()

	p2pOpts := []libp2p.Option{
		libp2p.EnableRelayService(),
	}

	p, err = p2p.Init(nil, p2pOpts...)
	if err != nil {
		log.Fatalf("p2p.Init: failed to initialize p2p: %v", err)
	}
	log.Info("libp2p node created: ", p.Node.ID())

	// Start a continous discovery process in the background
	// relay shouldn't require peers initially, so this'll just keep running.
	go p.DiscoveryLoop(ctx)

	http.HandleFunc("/", webHandler)

	log.Infof("Serving info on %s", config.GetHttpSocket())
	err = http.ListenAndServe(config.GetHttpSocket(), nil)
	if err != nil {
		log.Fatal(err)
	}

}
