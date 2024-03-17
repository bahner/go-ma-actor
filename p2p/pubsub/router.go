package pubsub

import (
	p2pubsub "github.com/libp2p/go-libp2p-pubsub"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

const PUBSUB_PROTOCOL = "/meshsub/1.1.0"

var r *pubsub.GossipSubRouter

func newRouter(h host.Host) *pubsub.GossipSubRouter {

	r = p2pubsub.DefaultGossipSubRouter(h)

	return r
}

// Adds a peer ID to the GossipSubRouter
func AddPeer(id peer.ID) {
	r.AddPeer(id, PUBSUB_PROTOCOL)
	r.AcceptFrom(id)

}

// Removes a peer ID from the GossipSubRouter
func RemovePeer(id peer.ID) {
	r.RemovePeer(id)
}

func SetEoughPeers(t string, peerno int) {
	r.EnoughPeers(t, peerno)
}
