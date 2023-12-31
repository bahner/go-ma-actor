package node

import (
	"fmt"
	"strconv"

	ipnskey "github.com/bahner/go-ma/key/ipns"
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
// Requires an IPNS key for identity libp2p options as parameters.
func New(i *ipnskey.Key, opts ...libp2p.Option) (host.Host, error) {

	p2pOptions := []libp2p.Option{
		libp2p.ListenAddrStrings(getListenAddrStrings()...),
		libp2p.Identity(i.PrivKey),
	}

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

	// Converting this to s string so quickly, is a little ugly,
	// but it is an integer to begin with, co it feels more correct.
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
