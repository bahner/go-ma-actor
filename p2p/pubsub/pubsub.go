package pubsub

import (
	"context"
	"fmt"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	log "github.com/sirupsen/logrus"
)

var pubSubService *pubsub.PubSub

func New(ctx context.Context, h host.Host) (*pubsub.PubSub, error) {

	pubSubService, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		return nil, fmt.Errorf("p2p: failed to create pubsub service: %v", err)
	}
	log.Info("Global resources initialized")

	return pubSubService, nil

}

func Get() (*pubsub.PubSub, error) {

	if pubSubService == nil {
		return nil, fmt.Errorf("p2p: pubsub service not initialized")
	}

	return pubSubService, nil
}
