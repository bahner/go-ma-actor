package p2p

import (
	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/db"
	libp2p "github.com/libp2p/go-libp2p"
	p2pDHT "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
)

type Options struct {
	DHT     []p2pDHT.Option
	P2P     []libp2p.Option
	Connmgr []connmgr.Option
}

// Default options for libp2p and DHT. This requires a private key,
// as we never create options without, so ...
func DefaultP2POptions() *Options {

	identity, err := db.GetOrCreateIdentity(config.ActorNick())
	if err != nil {
		panic(err)
	}

	return &Options{
		DHT: []p2pDHT.Option{
			p2pDHT.Mode(p2pDHT.ModeAutoServer),
		},
		P2P: []libp2p.Option{
			libp2p.Defaults,
			libp2p.Identity(identity),
		},
		Connmgr: []connmgr.Option{},
	}
}
