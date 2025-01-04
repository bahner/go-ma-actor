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
		libp2p.ListenAddrStrings(config.P2PMaddrs()...),
		libp2p.Identity(i),
	}

	log.Debugf("p2p: listen addresses: %v", config.P2PMaddrs())

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
