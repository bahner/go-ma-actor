package p2p

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/p2p/connmgr"
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
	PubSub   *p2ppubsub.PubSub
	DHT      *DHT
	MDNS     *MDNS
	Host     host.Host
	AddrInfo p2peer.AddrInfo
}

// Initialise everything needed for p2p communication. The function forces use of a specific IPNS key.
// Taken from the config package. It would be an error to initialise the node with a different key.
// The input is derived from Config() in the config package.
func Init(opts Options) (*P2P, error) {

	ctx := context.Background()

	// Initialise the connection manager and the connection gater.
	// Which needs to be passed to the libp2p node.
	cm, err := connmgr.Init(opts.Connmgr...)
	if err != nil {
		panic(fmt.Errorf("pong: failed to create connection manager: %w", err))
	}
	cg := connmgr.NewConnectionGater(cm)

	opts.P2P = append(opts.P2P,
		libp2p.ConnectionGater(cg),
	)

	// Initialise the libp2p node with the options.
	n, err := node.New(config.NodeIdentity(), opts.P2P...)
	if err != nil {
		return nil, fmt.Errorf("pong: failed to create libp2p node: %w", err)
	}

	err = initPeer(string(n.ID()))
	if err != nil {
		return nil, fmt.Errorf("p2p.Init: failed to init peer: %w", err)
	}

	// Create discovery and routing
	d, err := NewDHT(n, cg)
	if err != nil {
		return nil, fmt.Errorf("pong: failed to create DHT: %w", err)
	}

	m, err := newMDNS(n)
	if err != nil {
		log.Errorf("p2p.Init: failed to start MDNS discovery: %v", err)
	}

	// Create the pubsub service
	ps, err := pubsub.New(ctx, d.Host)
	if err != nil {
		return nil, fmt.Errorf("p2p.Init: failed to create pubsub: %w", err)
	}

	// Generate the struct and start the protection loop.
	ai := p2peer.AddrInfo{
		ID:    d.Host.ID(),
		Addrs: d.Host.Addrs(),
	}

	_p2p = &P2P{
		AddrInfo: ai,
		DHT:      d,
		Host:     d.Host,
		MDNS:     m,
		PubSub:   ps,
	}

	go _p2p.protectLoop(ctx)

	return _p2p, nil
}

// Get the P2P instance. This is a singleton.
func Get() *P2P {
	return _p2p
}

func (p *P2P) StartDiscoveryLoop(ctx context.Context) error {

	if config.P2PDiscoveryMDNS() {
		p.MDNS.discoveryLoop(ctx)
	}

	if config.P2PDiscoveryDHT() {
		go p.DHT.discoveryLoop(ctx)
	}

	return nil
}
