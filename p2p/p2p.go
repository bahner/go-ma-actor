package p2p

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/p2p/connmgr"
	"github.com/bahner/go-ma-actor/p2p/dht"
	"github.com/bahner/go-ma-actor/p2p/node"
	"github.com/bahner/go-ma-actor/p2p/pubsub"
	libp2p "github.com/libp2p/go-libp2p"
	p2ppubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	p2peer "github.com/libp2p/go-libp2p/core/peer"
	log "github.com/sirupsen/logrus"
)

var _p2p *P2P

// This is not a normal libp2p node, it's a wrapper around it. And it is specific to this project.
// It contains a libp2p node, a pubsub service and a DHT instance.
// It also contains a list of connected peers.
type P2P struct {
	Node     host.Host
	PubSub   *p2ppubsub.PubSub
	DHT      *dht.DHT
	AddrInfo *p2peer.AddrInfo
}

// Initialise everything needed for p2p communication. The function forces use of a specific IPNS key.
// Taken from the config package. It would be an error to initialise the node with a different key.
//
// d is a DHT instance. If nil, a new one will be created.
//
// Also takes a variadic list of libp2p options.
// Of it's nil, an empty list will be used.
//
// The configurable connection manager will be added to the node.
//
// The function return the libp2p node and a PubSub Service

func Init(d *dht.DHT, p2pOpts ...libp2p.Option) (*P2P, error) {

	log.Debug("p2p: Initialising p2p")

	var err error

	cm, err := connmgr.Init()
	if err != nil {
		return nil, fmt.Errorf("p2p.Init: failed to create connection manager: %w", err)
	}
	p2pOpts = append(p2pOpts, libp2p.ConnectionManager(cm))

	cg := connmgr.NewConnectionGater(cm)
	p2pOpts = append(p2pOpts, libp2p.ConnectionGater(cg))

	log.Debug("p2p: Connection manager initialised as follows:")
	log.Debug(cm.GetInfo())

	// Create a new libp2p Host that listens on a random TCP port
	n, err := node.New(config.NodeIdentity(), p2pOpts...)
	if err != nil {
		return nil, fmt.Errorf("p2p.Init: failed to create libp2p node: %w", err)
	}

	if d == nil {

		d, err = dht.New(n, cg)
		if err != nil {
			return nil, fmt.Errorf("p2p.Init: failed to create DHT: %w", err)
		}
	}

	ps, err := pubsub.New(context.Background(), n)
	if err != nil {
		return nil, fmt.Errorf("p2p.Init: failed to create pubsub: %w", err)
	}

	ai := p2peer.AddrInfo{
		ID:    n.ID(),
		Addrs: n.Addrs(),
	}

	_p2p = &P2P{
		AddrInfo: &ai,
		DHT:      d,
		Node:     n,
		PubSub:   ps,
	}

	return _p2p, nil
}

// Get the P2P instance. This is a singleton.
func Get() *P2P {
	return _p2p
}
