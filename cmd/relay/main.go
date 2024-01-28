package main

import (
	"context"
	"net/http"
	"time"

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

	// Boostrap Kademlia DHT and wait for it to finish.
	go discoveryLoop(ctx, p)

	http.HandleFunc("/", webHandler)

	log.Infof("Serving info on %s", getHttpSocket())
	err = http.ListenAndServe(getHttpSocket(), nil)
	if err != nil {
		log.Fatal(err)
	}

}

func discoveryLoop(ctx context.Context, p *p2p.P2P) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			p.DHT.DiscoverPeers()
			sleepTime := config.GetDiscoveryRetryInterval()
			log.Debugf("Sleeping for %v", sleepTime)
			log.Debugf(config.GetDiscoveryRetryIntervalString())
			time.Sleep(sleepTime)
		}
	}
}
