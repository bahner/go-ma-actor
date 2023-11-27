package config

import (
	"context"
	"flag"
	"time"

	"go.deanishe.net/env"
)

var (
	// P2P Settings
	lowWaterMark = flag.Int("lowWaterMark", env.GetInt(GO_MA_ACTOR_LOW_WATERMARK_VAR, defaultLowWaterMark),
		"Low watermark for peer discovery. You can use environment variable "+GO_MA_ACTOR_LOW_WATERMARK_VAR+" to set this.")
	highWaterMark = flag.Int("highWaterMark", env.GetInt(GO_MA_ACTOR_HIGH_WATERMARK_VAR, defaultHighWaterMark),
		"High watermark for peer discovery. You can use environment variable "+GO_MA_ACTOR_HIGH_WATERMARK_VAR+" to set this.")
	desiredPeers = flag.Int("desiredPeers", env.GetInt(GO_MA_ACTOR_DESIRED_PEERS_VAR, defaultDesiredPeers),
		"Desired number of peers to connect to. You can use environment variable "+GO_MA_ACTOR_DESIRED_PEERS_VAR+" to set this.")

	connmgrGracePeriod = flag.Duration("connmgrGracePeriod", env.GetDuration(GO_MA_ACTOR_CONNMGR_GRACE_VAR, defaultConnMgrGrace),
		"Grace period for connection manager. You can use environment variable "+GO_MA_ACTOR_CONNMGR_GRACE_VAR+" to set this.")
	discoveryRetryInterval = flag.Duration("discoveryRetryInterval", env.GetDuration(GO_MA_ACTOR_DISCOVERY_RETRY_INTERVAL_VAR, defaultDiscoveryRetryInterval),
		"Retry interval for peer discovery. You can use environment variable "+GO_MA_ACTOR_DISCOVERY_RETRY_INTERVAL_VAR+" to set this.")
	discoveryTimeout = flag.Duration("discoveryTimeout", env.GetDuration(GO_MA_ACTOR_DISCOVERY_TIMEOUT_VAR, defaultDiscoveryTimeout),
		"Timeout for peer discovery. You can use environment variable "+GO_MA_ACTOR_DISCOVERY_TIMEOUT_VAR+" to set this.")
)

// P2P Settings
func GetDiscoveryTimeout() time.Duration {
	return time.Duration(*discoveryTimeout) * time.Second
}

func GetLowWaterMark() int {
	return *lowWaterMark
}

func GetHighWaterMark() int {
	return *highWaterMark
}

func GetConnMgrGracePeriod() time.Duration {
	return *connmgrGracePeriod
}

func GetDiscoveryContext() (context.Context, func()) {

	ctx := context.Background()

	discoveryCtx, cancel := context.WithTimeout(ctx, GetDiscoveryTimeout())

	return discoveryCtx, cancel
}

func GetDiscoveryRetryInterval() time.Duration {
	return *discoveryRetryInterval
}

func GetDesiredPeers() int {
	return *desiredPeers
}
