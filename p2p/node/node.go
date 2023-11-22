package node

import (
	"fmt"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
)

var p2pNode host.Host

// Creates a new libp2p node, meant to be the only one used in an application.
// Takes normal libp2p options as parameters.
func New(opts ...libp2p.Option) (host.Host, error) {

	// Create a new libp2p Host that listens on a random TCP port
	p2pNode, err := libp2p.New(opts...)

	if err != nil {
		return nil, fmt.Errorf("p2p: failed to create libp2p node: %v", err)
	}

	return p2pNode, nil
}

func Get() (host.Host, error) {

	if p2pNode == nil {
		return nil, fmt.Errorf("p2p: node not initialized")
	}

	return p2pNode, nil
}
