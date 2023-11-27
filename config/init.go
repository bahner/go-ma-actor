package config

import (
	"context"
	"flag"
	"os"
	"time"

	"github.com/bahner/go-ma/key/set"
	log "github.com/sirupsen/logrus"
)

var (
	name = "go-ma-actor"

	keyset *set.Keyset

	generate     *bool
	genenv       *bool
	publish      *bool
	forcePublish *bool
)

func init() {

	// Booleans with control flow
	generate = flag.Bool("generate", false, "Generates one-time keyset and uses it")
	genenv = flag.Bool("genenv", false, "Generates a keyset and prints it to stdout and uses it")
	publish = flag.Bool("publish", false, "Publishes keyset to IPFS when using genenv or generate")
	forcePublish = flag.Bool("forcePublish", false, "Like -publish, force publication even if keyset is already published. This is probably the one you want.")

	// Entities
	flag.StringVar(&nick, "nick", nick, "Nickname to use in character creation")
	flag.StringVar(&keyset_string, "keyset", keyset_string, "Base58 encoded *secret* keyset used to identify the client. You. You can use environment variable "+GO_MA_ACTOR_KEYSET_VAR+" to set this.")
	flag.StringVar(&entity, "entity", entity, "DID of the entity to communicate with. You can use environment variable "+GO_MA_ACTOR_ENTITY_VAR+" to set this.")

	// P2P Settings
	flag.IntVar(&lowWaterMark, "lowWaterMark", lowWaterMark, "Low watermark for peer discovery. You can use environment variable "+GO_MA_ACTOR_LOW_WATERMARK_VAR+" to set this.")
	flag.IntVar(&highWaterMark, "highWaterMark", highWaterMark, "High watermark for peer discovery. You can use environment variable "+GO_MA_ACTOR_HIGH_WATERMARK_VAR+" to set this.")
	flag.IntVar(&desiredPeers, "desiredPeers", desiredPeers, "Desired number of peers to connect to. You can use environment variable "+GO_MA_ACTOR_DESIRED_PEERS_VAR+" to set this.")

	flag.DurationVar(&connmgrGracePeriod, "connmgrGracePeriod", connmgrGracePeriod, "Grace period for connection manager. You can use environment variable "+GO_MA_ACTOR_CONNMGR_GRACE_VAR+" to set this.")
	flag.DurationVar(&discoveryRetryInterval, "discoveryRetryInterval", discoveryRetryInterval, "Retry interval for peer discovery. You can use environment variable "+GO_MA_ACTOR_DISCOVERY_RETRY_INTERVAL_VAR+" to set this.")
	flag.DurationVar(&discoveryTimeout, "discoveryTimeout", discoveryTimeout, "Timeout for peer discovery. You can use environment variable "+GO_MA_ACTOR_DISCOVERY_TIMEOUT_VAR+" to set this.")

	// Logging
	flag.StringVar(&logLevel, "loglevel", logLevel, "Loglevel to use for application. You can use environment variable "+GO_MA_ACTOR_LOGLEVEL_VAR+" to set this.")
	flag.StringVar(&logfile, "logfile", logfile, "Logfile to use for application. You can use environment variable "+GO_MA_ACTOR_LOGFILE_VAR+" to set this.")

	flag.Parse()

	// Init logger
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Error(err)
		os.Exit(64) // EX_USAGE
	}
	log.SetLevel(level)
	file, err := os.OpenFile(name+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Errorf("Failed to open log file: %v", err)
		os.Exit(73) // EX_CANTCREAT
	}
	log.SetOutput(file)

	log.Info("Logger initialized")

	// Init keyset
	initKeyset(keyset_string)

}

func GetNick() string {
	return nick
}

func GetEntity() string {
	return entity
}

func GetLogLevel() string {
	return logLevel
}

func GetPublish() bool {

	return *publish
}

func GetForcePublish() bool {
	return *forcePublish
}

// P2P Settings
func GetDiscoveryTimeout() time.Duration {
	return time.Duration(discoveryTimeout) * time.Second
}

func GetLowWaterMark() int {
	return lowWaterMark
}

func GetHighWaterMark() int {
	return highWaterMark
}

func GetConnMgrGracePeriod() time.Duration {
	return connmgrGracePeriod
}

func GetDiscoveryContext() (context.Context, func()) {

	ctx := context.Background()

	discoveryCtx, cancel := context.WithTimeout(ctx, GetDiscoveryTimeout())

	return discoveryCtx, cancel
}

func GetDiscoveryRetryInterval() time.Duration {
	return discoveryRetryInterval
}

func GetDesiredPeers() int {
	return desiredPeers
}
