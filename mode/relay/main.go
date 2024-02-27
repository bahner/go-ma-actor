package relay

import (
	libp2p "github.com/libp2p/go-libp2p"
)

// Just enable relay for now. I expect much more logic here in the future.
func GetOptions() []libp2p.Option {
	p2pOpts := []libp2p.Option{
		libp2p.EnableRelayService(),
	}

	return p2pOpts
}
