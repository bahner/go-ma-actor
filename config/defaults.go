package config

import (
	"time"

	"github.com/bahner/go-ma"
	"github.com/bahner/go-ma/key/set"
	"go.deanishe.net/env"
)

const (
	name = "go-ma-actor"

	// The default entity to connect to.
	GO_MA_ACTOR_KEYSET_VAR            = "GO_MA_ACTOR_KEYSET"
	GO_MA_ACTOR_ENTITY_VAR            = "GO_MA_ACTOR_ENTITY"
	GO_MA_ACTOR_DISCOVERY_TIMEOUT_VAR = "GO_MA_ACTOR_DISCOVERY_TIMEOUT"
	GO_MA_ACTOR_LOW_WATERMARK_VAR     = "GO_MA_ACTOR_LOW_WATERMARK"
	GO_MA_ACTOR_HIGH_WATERMARK_VAR    = "GO_MA_ACTOR_HIGH_WATERMARK"
	GO_MA_ACTOR_CONNMGR_GRACE_VAR     = "GO_MA_ACTOR_CONNMGR_GRACE"

	defaultDiscoveryTimeout int           = 300
	defaultLowWaterMark     int           = 2
	defaultHighWaterMark    int           = 10
	defaultConnMgrGrace     time.Duration = time.Minute * 1
)

var (
	generate     *bool
	genenv       *bool
	publish      *bool
	forcePublish *bool

	keyset *set.Keyset

	discoveryTimeout   int           = env.GetInt(GO_MA_ACTOR_DISCOVERY_TIMEOUT_VAR, defaultDiscoveryTimeout)
	lowWaterMark       int           = env.GetInt(GO_MA_ACTOR_LOW_WATERMARK_VAR, defaultLowWaterMark)
	highWaterMark      int           = env.GetInt(GO_MA_ACTOR_HIGH_WATERMARK_VAR, defaultHighWaterMark)
	connmgrGracePeriod time.Duration = env.GetDuration(GO_MA_ACTOR_CONNMGR_GRACE_VAR, defaultConnMgrGrace)

	logLevel string = env.Get(ma.LOGLEVEL_VAR, "info")
	logfile  string = env.Get(ma.LOGFILE_VAR, name+".log")

	// What we want to communicate with initially
	entity string = env.Get(GO_MA_ACTOR_ENTITY_VAR, "")

	// Actor
	keyset_string string = env.Get(GO_MA_ACTOR_KEYSET_VAR, "")

	// Nick is only used for keyset generation. Must be a valid NanoID.
	nick string = env.Get("USER")
)
