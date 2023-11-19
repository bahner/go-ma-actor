package p2p

import (
	"context"
	"fmt"

	"github.com/bahner/go-home/config"
	"github.com/bahner/go-ma/key/set"
	"github.com/bahner/go-space/p2p/host"
	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

func Init(ctx context.Context, k *set.Keyset) (*host.P2pHost, *pubsub.PubSub, error) {

	node, err := host.New(
		libp2p.Identity(k.IPNSKey.PrivKey),
		libp2p.ListenAddrStrings(),
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to create p2p host: %v", err))
	}

	node.StartPeerDiscovery(ctx, config.GetRendezvous())

	ps, err := pubsub.NewGossipSub(ctx, node)
	if err != nil {
		panic(fmt.Sprintf("Failed to create pubsub: %v", err))
	}

	return node, ps, nil

}
