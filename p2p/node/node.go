package node

import (
	"strconv"

	"github.com/bahner/go-ma-actor/config"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	log "github.com/sirupsen/logrus"
)

const NODE_LISTEN_PORT = 4001

var (
	err     error
	p2pNode host.Host
)

// Creates a new libp2p node, meant to be the only one used in an application.
// Takes normal libp2p options as parameters.
func init() {

	// Get the keyset and use it to start the node
	k := config.GetKeyset()

	// Create a new libp2p Host that listens on a random TCP port
	p2pNode, err = libp2p.New(
		libp2p.ListenAddrStrings(getListenAddrStrings()...),
		libp2p.Identity(k.IPNSKey.PrivKey),
	)

	if err != nil {
		log.Errorf("p2p: failed to create libp2p node: %v", err)
	}
}

func Get() host.Host {

	if p2pNode == nil {
		log.Errorf("p2p: node not initialized")
	}

	return p2pNode
}

func getListenAddrStrings() []string {

	port := strconv.Itoa(NODE_LISTEN_PORT)

	return []string{
		"/ip4/0.0.0.0/tcp/" + port,
		"/ip4/0.0.0.0/udp/" + port + "/quic",
		"/ip4/0.0.0.0/udp/" + port + "/quic-v1",
		"/ip4/0.0.0.0/udp/" + port + "/quic-v1/webtransport",
		"/ip6/::/tcp/" + port,
		"/ip6/::/udp/" + port + "/quic",
		"/ip6/::/udp/" + port + "/quic-v1",
		"/ip6/::/udp/" + port + "/quic-v1/webtransport",
	}
}
