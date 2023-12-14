package config

import (
	"flag"
	"fmt"
)

const (

	// Environment variable names
	GO_MA_ACTOR_CONNMGR_GRACE_VAR            = "GO_MA_ACTOR_CONNMGR_GRACE"
	GO_MA_ACTOR_DESIRED_PEERS_VAR            = "GO_MA_ACTOR_DESIRED_PEERS"
	GO_MA_ACTOR_DISCOVERY_RETRY_INTERVAL_VAR = "GO_MA_ACTOR_DISCOVERY_RETRY"
	GO_MA_ACTOR_DISCOVERY_TIMEOUT_VAR        = "GO_MA_ACTOR_DISCOVERY_TIMEOUT"
	GO_MA_ACTOR_ENTITY_VAR                   = "GO_MA_ACTOR_ENTITY"
	GO_MA_ACTOR_HIGH_WATERMARK_VAR           = "GO_MA_ACTOR_HIGH_WATERMARK"
	GO_MA_ACTOR_IDENTITY_VAR                 = "GO_MA_ACTOR_IDENTITY"
	GO_MA_ACTOR_LOGFILE_VAR                  = "GOLOG_FILE"
	GO_MA_ACTOR_LOGLEVEL_VAR                 = "GOLOG_LOG_LEVEL"
	GO_MA_ACTOR_LOW_WATERMARK_VAR            = "GO_MA_ACTOR_LOW_WATERMARK"
	GO_MA_ACTOR_NODE_IDENTITY_VAR            = "GO_MA_ACTOR_NODE_IDENTITY"
)

var genenv = flag.Bool("genenv", false, "Generates and prints environment to stdout")

func PrintEnvironment() {
	fmt.Println("# Entity Settings - the default actor to connect with")
	fmt.Println("export " + GO_MA_ACTOR_ENTITY_VAR + "=" + GetEntity())
	fmt.Println("# Secrets of the actor")
	fmt.Println("export " + GO_MA_ACTOR_IDENTITY_VAR + "=" + GetIdentityString())
	fmt.Println("# identity of the local node")
	fmt.Println("export " + GO_MA_ACTOR_NODE_IDENTITY_VAR + "=" + GetNodeMultibasePrivKey())
	fmt.Println("# P2P Settings")
	fmt.Println("export " + GO_MA_ACTOR_CONNMGR_GRACE_VAR + "=" + GetConnMgrGraceString())
	fmt.Println("export " + GO_MA_ACTOR_LOW_WATERMARK_VAR + "=" + GetLowWatermarkString())
	fmt.Println("export " + GO_MA_ACTOR_HIGH_WATERMARK_VAR + "=" + GetHighWatermarkString())
	fmt.Println("export " + GO_MA_ACTOR_DESIRED_PEERS_VAR + "=" + GetDesiredPeersString())
	fmt.Println("export " + GO_MA_ACTOR_DISCOVERY_RETRY_INTERVAL_VAR + "=" + GetDiscoveryRetryIntervalString())
	fmt.Println("export " + GO_MA_ACTOR_DISCOVERY_TIMEOUT_VAR + "=" + GetDiscoveryTimeoutString())
	fmt.Println("# Logging Settings")
	fmt.Println("export " + GO_MA_ACTOR_LOGFILE_VAR + "=" + GetLogFile())
	fmt.Println("export " + GO_MA_ACTOR_LOGLEVEL_VAR + "=" + GetLogLevel())
}
