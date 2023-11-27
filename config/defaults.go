package config

import (
	"time"
)

const (

	// The name of the application
	name = "go-ma-actor"

	// The default entity to connect to.
	GO_MA_ACTOR_KEYSET_VAR                   = "GO_MA_ACTOR_KEYSET"
	GO_MA_ACTOR_ENTITY_VAR                   = "GO_MA_ACTOR_ENTITY"
	GO_MA_ACTOR_DISCOVERY_TIMEOUT_VAR        = "GO_MA_ACTOR_DISCOVERY_TIMEOUT"
	GO_MA_ACTOR_DISCOVERY_RETRY_INTERVAL_VAR = "GO_MA_ACTOR_DISCOVERY_RETRY"
	GO_MA_ACTOR_DESIRED_PEERS_VAR            = "GO_MA_ACTOR_DESIRED_PEERS"
	GO_MA_ACTOR_LOW_WATERMARK_VAR            = "GO_MA_ACTOR_LOW_WATERMARK"
	GO_MA_ACTOR_HIGH_WATERMARK_VAR           = "GO_MA_ACTOR_HIGH_WATERMARK"
	GO_MA_ACTOR_CONNMGR_GRACE_VAR            = "GO_MA_ACTOR_CONNMGR_GRACE"
	GO_MA_ACTOR_LOGLEVEL_VAR                 = "GO_MA_ACTOR_LOGLEVEL"
	GO_MA_ACTOR_LOGFILE_VAR                  = "GO_MA_ACTOR_LOGFILE"

	defaultLowWaterMark  int = 3
	defaultHighWaterMark int = 10
	defaultDesiredPeers  int = 3

	defaultDiscoveryTimeout       time.Duration = time.Second * 30
	defaultConnMgrGrace           time.Duration = time.Minute * 1
	defaultDiscoveryRetryInterval time.Duration = time.Second * 1

	defaultLogLevel string = "info"
	defaultLogfile  string = name + ".log"

	defaultNick   string = "ghost"
	defaultKeyset string = ""
	defaultEntity string = ""
)
