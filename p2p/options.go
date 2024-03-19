package p2p

import (
	libp2p "github.com/libp2p/go-libp2p"
	p2pDHT "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
)

type Options struct {
	DHT     []p2pDHT.Option
	P2P     []libp2p.Option
	Connmgr []connmgr.Option
}

func DefaultOptions() Options {
	return Options{
		DHT: []p2pDHT.Option{
			p2pDHT.Mode(p2pDHT.ModeAutoServer),
		},
		P2P: []libp2p.Option{
			libp2p.DefaultTransports,
			libp2p.DefaultSecurity,
		},
		Connmgr: []connmgr.Option{},
	}
}
