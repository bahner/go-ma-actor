package config

import (
	"time"
)

const (
	defaultLowWaterMark  int = 3
	defaultHighWaterMark int = 10
	defaultDesiredPeers  int = 3

	defaultDiscoveryTimeout       time.Duration = time.Second * 30
	defaultConnMgrGrace           time.Duration = time.Minute * 1
	defaultDiscoveryRetryInterval time.Duration = time.Second * 1

	defaultLogLevel string = "info"
	defaultLogfile  string = Name + ".log"

	defaultNick   string = "ghost"
	defaultKeyset string = ""
	defaultEntity string = ""
)
