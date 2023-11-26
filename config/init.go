package config

import (
	"flag"
	"os"
	"time"

	"github.com/bahner/go-ma"
	"github.com/bahner/go-ma/key/set"
	log "github.com/sirupsen/logrus"
	"go.deanishe.net/env"
)

const (
	name = "go-ma-actor"

	// The default entity to connect to.
	GO_MA_ACTOR_ENTITY_VAR = "GO_MA_ACTOR_ENTITY"

	defaultDiscoveryTimeout int           = 300
	defaultLowWaterMark     int           = 2
	defaultHighWaterMark    int           = 10
	defaultConnMgrGrace     time.Duration = time.Minute * 1
)

var (
	err error

	generate     *bool
	genenv       *bool
	publish      *bool
	forcePublish *bool

	keyset *set.Keyset
)

var (
	discoveryTimeout   int           = env.GetInt(ma.DISCOVERY_TIMEOUT_VAR, defaultDiscoveryTimeout)
	lowWaterMark       int           = env.GetInt(ma.LOW_WATERMARK_VAR, defaultLowWaterMark)
	highWaterMark      int           = env.GetInt(ma.HIGH_WATERMARK_VAR, defaultHighWaterMark)
	connmgrGracePeriod time.Duration = env.GetDuration(ma.CONNMGR_GRACE_VAR, defaultConnMgrGrace)

	logLevel string = env.Get(ma.LOGLEVEL_VAR, "info")
	logfile  string = env.Get(ma.LOGFILE_VAR, name+".log")

	// What we want to communicate with initially
	entity string = env.Get(GO_MA_ACTOR_ENTITY_VAR, "")

	// Actor
	keyset_string string = env.Get(ma.KEYSET_VAR, "")

	// Nick is only used for keyset generation. Must be a valid NanoID.
	nick string = env.Get("USER")
)

func init() {

	// Flags - user configurations
	flag.StringVar(&logLevel, "loglevel", logLevel, "Loglevel to use for application")
	flag.StringVar(&logfile, "logfile", logfile, "Logfile to use for application")

	// P"P Settings
	flag.IntVar(&discoveryTimeout, "discoveryTimeout", discoveryTimeout, "Timeout for peer discovery")
	flag.IntVar(&lowWaterMark, "lowWaterMark", lowWaterMark, "Low watermark for peer discovery")
	flag.IntVar(&highWaterMark, "highWaterMark", highWaterMark, "High watermark for peer discovery")
	flag.DurationVar(&connmgrGracePeriod, "connmgrGracePeriod", connmgrGracePeriod, "Grace period for connection manager")

	// Actor
	flag.StringVar(&nick, "nick", nick, "Nickname to use in character creation")
	flag.StringVar(&keyset_string, "keyset", keyset_string, "Base58 encoded *secret* keyset used to identify the client. You.")
	flag.StringVar(&entity, "entity", entity, "DID of the entity to communicate with.")

	// Booleans with control flow
	generate = flag.Bool("generate", false, "Generates one-time keyset and uses it")
	genenv = flag.Bool("genenv", false, "Generates a keyset and prints it to stdout and uses it")
	publish = flag.Bool("publish", false, "Publishes keyset to IPFS when using genenv or generate")
	forcePublish = flag.Bool("forcePublish", false, "Like -publish, force publication even if keyset is already published. This is probably the one you want.")

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
