package config

import (
	"flag"
	"os"
	"time"

	"github.com/bahner/go-ma/key/set"
	log "github.com/sirupsen/logrus"
	"go.deanishe.net/env"
)

const (
	name                    = "go-ma-actor"
	keyset_var              = "GO_ACTOR_KEYSET"
	entity_var              = "GO_ACTOR_ENTITY"
	discovery_timeout_var   = "GO_ACTOR_DISCOVERY_TIMEOUT"
	log_level_var           = "GO_ACTOR_LOG_LEVEL"
	defaultDiscoveryTimeout = 300
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
	discoveryTimeout int    = env.GetInt(discovery_timeout_var, defaultDiscoveryTimeout)
	logLevel         string = env.Get(log_level_var, "info")
	logfile          string = env.Get("GO_ACTOR_LOG_FILE", name+"log")

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
	flag.StringVar(&logfile, "logfile", logfile, "Logfile to use for application")
	flag.IntVar(&discoveryTimeout, "discoveryTimeout", discoveryTimeout, "Timeout for peer discovery")

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
		log.Fatal(err)
	}
	log.SetLevel(level)
	file, err := os.OpenFile(name+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
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
