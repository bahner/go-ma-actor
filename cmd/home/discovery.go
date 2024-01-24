package main

import (
	"time"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/p2p"
)

func discoveryHandler(p *p2p.P2P) {
	for {
		p.DiscoverPeers()
		time.Sleep(config.GetDiscoveryTimeout())
	}

}
