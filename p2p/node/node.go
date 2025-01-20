package node

import (
	"fmt"

	"github.com/bahner/go-ma-actor/config"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	log "github.com/sirupsen/logrus"
)

var (
	err     error
	p2pNode host.Host
)

// Creates a new libp2p node, meant to be the only one used in an application.
// Requires an IPNS key for identity libp2p options as parameters.
func New(opts ...libp2p.Option) (host.Host, error) {

	log.Debugf("p2p: listen addresses: %v", config.P2PMaddrs())

	// Create a new libp2p Host with provided options
	p2pNode, err = libp2p.New(opts...)

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
