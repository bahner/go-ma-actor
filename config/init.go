package config

import (
	"flag"

	"github.com/bahner/go-ma/key/set"
	log "github.com/sirupsen/logrus"
	"go.deanishe.net/env"
)

const (
	keyset_var              = "GO_ACTOR_KEYSET"
	entity_var              = "GO_ACTOR_ENTITY"
	discovery_timeout_var   = "GO_ACTOR_DISCOVERY_TIMEOUT"
	log_level_var           = "GO_ACTOR_LOG_LEVEL"
	defaultDiscoveryTimeout = 300
)

var (
	generate     *bool
	genenv       *bool
	publish      *bool
	forcePublish *bool

	keyset *set.Keyset
)

var (
	discoveryTimeout int    = env.GetInt(discovery_timeout_var, defaultDiscoveryTimeout)
	logLevel         string = env.Get(log_level_var, "error")

	// What we want to communicate with initially
	entity string = env.Get(entity_var, "")

	// Actor
	keyset_string string = env.Get(keyset_var, "")

	// Nick is only used for keyset generation. Must be a valid NanoID.
	nick string = env.Get("USER")
)

func init() {

	// Flags - user configurations
	flag.StringVar(&logLevel, "loglevel", logLevel, "Loglevel to use for application")
	flag.IntVar(&discoveryTimeout, "discoveryTimeout", discoveryTimeout, "Timeout for peer discovery")

	// Actor
	flag.StringVar(&nick, "nick", nick, "Nickname to use in character creation")
	flag.StringVar(&keyset_string, "keyset", keyset_string, "Base58 encoded secret key used to identify the client. You.")
	flag.StringVar(&entity, "entity", entity, "DID of the entity to communicate with.")

	// Booleans with control flow
	generate = flag.Bool("generate", false, "Generates one-time keyset and uses it")
	genenv = flag.Bool("genenv", false, "Generates a keyset and prints it to stdout and uses it")
	publish = flag.Bool("publish", false, "Publishes keyset to IPFS when using genenv or generate")
	forcePublish = flag.Bool("forcePublish", false, "Force publish even if keyset is already published")

	flag.Parse()

	// Init logger
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(level)
	log.Info("Logger initialized")

	// Init keyset
	initKeyset(keyset_string)

	// Make sure required services are running
	initP2P(discoveryTimeout)

	// Init actor and libP2P node
	initActor() // Requires running IPFS Daemon for publication

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
