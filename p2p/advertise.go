package p2p

import (
	"context"
	"time"

	"github.com/bahner/go-ma"
	"github.com/bahner/go-ma-actor/config"
	"github.com/libp2p/go-libp2p/core/discovery"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	"github.com/libp2p/go-libp2p/p2p/discovery/util"
	log "github.com/sirupsen/logrus"
)

func advertisementLoop(ctx context.Context, routingDiscovery *routing.RoutingDiscovery, discoveryOpts ...discovery.Option) {

	ticker := time.NewTicker(config.P2PDiscoveryAdvertiseInterval())
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			util.Advertise(ctx, routingDiscovery, ma.RENDEZVOUS, discoveryOpts...)
			log.Debugf("Advertising rendezvous string: %s", ma.RENDEZVOUS)
		case <-ctx.Done():
			return
		}
	}
}
