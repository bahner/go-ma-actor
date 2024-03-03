package connmgr

import (
	"github.com/bahner/go-ma-actor/config"
	p2pConnmgr "github.com/libp2p/go-libp2p/p2p/net/connmgr"
	log "github.com/sirupsen/logrus"
)

var connmgr *p2pConnmgr.BasicConnMgr

func Init(opts ...p2pConnmgr.Option) (*p2pConnmgr.BasicConnMgr, error) {

	if connmgr != nil {
		log.Debugf("Connection manager already initialized")
		return connmgr, nil
	}

	lowWaterMark := config.P2PConnmgrLowWatermark()
	log.Infof("Connection manager low water mark: %v", lowWaterMark)

	highWaterMark := config.P2PConnmgrHighWatermark()
	log.Infof("Connection manager high water mark: %v", highWaterMark)

	gracePeriod := config.P2PConnMgrGracePeriod()
	opts = append(opts, p2pConnmgr.WithGracePeriod(gracePeriod))
	log.Infof("Connection manager grace period: %v", gracePeriod)

	return p2pConnmgr.NewConnManager(
		lowWaterMark,
		highWaterMark,
		opts...,
	)

}
