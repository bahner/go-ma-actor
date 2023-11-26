package connmgr

import (
	"github.com/bahner/go-ma-actor/config"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	p2pConnmgr "github.com/libp2p/go-libp2p/p2p/net/connmgr"
	log "github.com/sirupsen/logrus"
)

func Init() (*p2pConnmgr.BasicConnMgr, error) {

	gracePeriod := config.GetConnMgrGracePeriod()
	withGracePeriod := connmgr.WithGracePeriod(gracePeriod)
	log.Infof("Connection manager grace period: %v", gracePeriod)

	lowWaterMark := config.GetLowWaterMark()
	log.Infof("Connection manager low water mark: %v", lowWaterMark)
	highWaterMark := config.GetHighWaterMark()
	log.Infof("Connection manager high water mark: %v", highWaterMark)

	return p2pConnmgr.NewConnManager(
		lowWaterMark,
		highWaterMark,
		withGracePeriod,
	)

}
