package pubsub

import (
	"context"

	"github.com/bahner/go-ma-actor/p2p/node"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	log "github.com/sirupsen/logrus"
)

var (
	err error

	pubSubService *pubsub.PubSub
)

func init() {
	ctx := context.Background()
	n := node.Get()

	pubSubService, err = pubsub.NewGossipSub(ctx, n)
	if err != nil {
		log.Fatalf("p2p: failed to create pubsub service: %v", err)
	}
	log.Info("Global resources initialized")

}

func Get() *pubsub.PubSub {

	if pubSubService == nil {
		log.Errorf("p2p: pubsub service not initialized")
	}

	return pubSubService
}
