package pubsub

import (
	"context"
	"fmt"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	log "github.com/sirupsen/logrus"
)

var (
	err error

	ps *pubsub.PubSub
)

// Start the pubsub service. If ctx is nil a background context is used.
// n is the libp2p host.
func New(ctx context.Context, n host.Host) (*pubsub.PubSub, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	ps, err = pubsub.NewGossipSub(ctx, n)
	if err != nil {
		return nil, fmt.Errorf("p2p: failed to create pubsub service: %w", err)
	}

	return ps, nil
}

func Get() *pubsub.PubSub {

	if ps == nil {
		log.Errorf("p2p: pubsub service not initialized")
	}

	return ps
}
