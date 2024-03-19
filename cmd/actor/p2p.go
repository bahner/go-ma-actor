package p2p

import (
	"fmt"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/p2p"
	"github.com/bahner/go-ma-actor/p2p/connmgr"
	"github.com/bahner/go-ma-actor/p2p/node"
	"github.com/bahner/go-ma-actor/p2p/peer"
	"github.com/libp2p/go-libp2p"
)

func actorDHT(cg *connmgr.ConnectionGater) (*DHT, error) {

	// THese are the relay specific parts.
	p2pOpts := []libp2p.Option{
		libp2p.ConnectionGater(cg),
	}

	n, err := node.New(config.NodeIdentity(), p2pOpts...)
	if err != nil {
		return nil, fmt.Errorf("pong: failed to create libp2p node: %w", err)
	}

	d, err := NewDHT(n, cg)
	if err != nil {
		return nil, fmt.Errorf("pong: failed to create DHT: %w", err)
	}

	return d, nil
}

func initP2P() (_p2p *P2P, err error) {
	fmt.Println("Initialising libp2p...")

	// Everyone needs a connection manager.
	cm, err := connmgr.Init()
	if err != nil {
		panic(fmt.Errorf("pong: failed to create connection manager: %w", err))
	}
	cg := connmgr.NewConnectionGater(cm)

	d, err := actorDHT(cg)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize dht: %v", err))
	}

	p, err := p2p.Init(d)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize p2p: %v", err))
	}

	// PEER
	fmt.Println("Initialising peer ...")
	err = initPeer(p.Host.ID().String())
	if err != nil {
		panic(fmt.Sprintf("failed to initialize peer: %v", err))
	}

	return p, nil

}

func initPeer(id string) error {

	p, err := peer.GetOrCreate(id)
	if err != nil {
		return fmt.Errorf("error getting or creating peer: %s", err)
	}

	return p.SetAllowed(true) // Ensure allowed even if we were previously denied.
}
