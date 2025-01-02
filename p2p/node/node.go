package node

import (
	"fmt"

	"github.com/bahner/go-ma-actor/config"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	log "github.com/sirupsen/logrus"
)

var (
	err     error
	p2pNode host.Host
)

// Creates a new libp2p node, meant to be the only one used in an application.
// Requires an IPNS key for identity libp2p options as parameters.
func New(i crypto.PrivKey, opts ...libp2p.Option) (host.Host, error) {

	p2pOptions := []libp2p.Option{
		libp2p.ListenAddrStrings(getListenAddrStrings()...),
		libp2p.Identity(i),
	}

	log.Debugf("p2p: listen addresses: %v", getListenAddrStrings())

	p2pOptions = append(p2pOptions, opts...)

	// Create a new libp2p Host that listens on a random TCP port
	p2pNode, err = libp2p.New(p2pOptions...)

	if err != nil {
		return nil, fmt.Errorf("p2p: failed to create libp2p node: %w", err)
	}

	return p2pNode, nil
}

func Get() host.Host {

	if p2pNode == nil {
		log.Errorf("p2p: node not initialized")
	}

	return p2pNode
}

func getListenAddrStrings() []string {

	port := config.P2PPortString()
	portws := config.P2PPortWSString()

	// This specifically adds "/quic" and "/ws" to the listen addresses.
	return []string{
		"/ip4/0.0.0.0/tcp/" + port,
		"/ip4/0.0.0.0/tcp/" + portws + "/ws",

		"/ip4/0.0.0.0/udp/" + port + "/quic",
		"/ip4/0.0.0.0/udp/" + port + "/quic-v1",
		"/ip4/0.0.0.0/udp/" + port + "/quic-v1/webtransport",

		"/ip6/::/tcp/" + port,
		"/ip6/::/tcp/" + portws + "/ws",

		"/ip6/::/udp/" + port + "/quic",
		"/ip6/::/udp/" + port + "/quic-v1",
		"/ip6/::/udp/" + port + "/quic-v1/webtransport",
	}
}
