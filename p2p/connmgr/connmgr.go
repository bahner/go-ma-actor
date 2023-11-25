package connmgr

import (
	"github.com/bahner/go-ma-actor/config"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	p2pConnmgr "github.com/libp2p/go-libp2p/p2p/net/connmgr"
)

func Init() (*p2pConnmgr.BasicConnMgr, error) {

	withGracePeriod := connmgr.WithGracePeriod(config.GetConnMgrGracePeriod())

	return p2pConnmgr.NewConnManager(
		config.GetLowWaterMark(),
		config.GetHighWaterMark(),
		withGracePeriod,
	)

}
