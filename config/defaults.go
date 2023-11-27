package config

import (
	"time"

	"github.com/bahner/go-ma"
	"go.deanishe.net/env"
)

const (

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

	defaultLowWaterMark  int = 3 // This will be used for the connection manager and the number of peers to search for
	defaultHighWaterMark int = 10
	defaultDesiredPeers  int = 3

	defaultDiscoveryTimeout       time.Duration = time.Second * 30
	defaultConnMgrGrace           time.Duration = time.Minute * 1
	defaultDiscoveryRetryInterval time.Duration = time.Second * 1
)

var (
	logLevel string = env.Get(ma.LOGLEVEL_VAR, "info")
	logfile  string = env.Get(ma.LOGFILE_VAR, name+".log")

	desiredPeers  int = env.GetInt(GO_MA_ACTOR_DESIRED_PEERS_VAR, defaultDesiredPeers)
	highWaterMark int = env.GetInt(GO_MA_ACTOR_HIGH_WATERMARK_VAR, defaultHighWaterMark)
	lowWaterMark  int = env.GetInt(GO_MA_ACTOR_LOW_WATERMARK_VAR, defaultLowWaterMark)

	connmgrGracePeriod     time.Duration = env.GetDuration(GO_MA_ACTOR_CONNMGR_GRACE_VAR, defaultConnMgrGrace)
	discoveryTimeout       time.Duration = env.GetDuration(GO_MA_ACTOR_DISCOVERY_TIMEOUT_VAR, defaultDiscoveryTimeout)
	discoveryRetryInterval time.Duration = env.GetDuration(GO_MA_ACTOR_DISCOVERY_RETRY_INTERVAL_VAR, defaultDiscoveryRetryInterval)

	// What we want to communicate with initially
	entity string = env.Get(GO_MA_ACTOR_ENTITY_VAR, "")

	// Actor
	keyset_string string = env.Get(GO_MA_ACTOR_KEYSET_VAR, "")

	// Nick is only used for keyset generation. Must be a valid NanoID.
	nick string = env.Get("USER")
)
