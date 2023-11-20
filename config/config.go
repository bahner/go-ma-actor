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
	keyset_var              = "GO_ACTOR_KEYSET"
	entity_var              = "GO_ACTOR_ENTITY"
	discovery_timeout_var   = "GO_ACTOR_DISCOVERY_TIMEOUT"
	log_level_var           = "GO_ACTOR_LOG_LEVEL"
	defaultDiscoveryTimeout = 300
)

var (
	discoveryTimeout int    = env.GetInt(discovery_timeout_var, defaultDiscoveryTimeout)
	logLevel         string = env.Get(log_level_var, "error")

	entity string = env.Get(entity_var, "")

	// Actor
	keyset_string string = env.Get(keyset_var, "")
	// Nick is only used for keyset generation. Must be a valid NanoID.
	nick string = env.Get("USER")

	generate     *bool
	genenv       *bool
	publish      *bool
	forcePublish *bool

	keyset *set.Keyset
)

func Init() {

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

	// Generate a new keysets if requested
	if *generate || *genenv {
		keyset_string = generateKeyset(keyset_var, nick, *forcePublish)
	}

	if *publish || *forcePublish {
		if keyset_string != "" {
			publishKeyset(keyset, *forcePublish)
		} else {
			log.Errorf("No actor keyset to publish.")
		}

		if *genenv || *generate {
			os.Exit(0)
		}

		log.Debugf("actor_keyset_string: %s", keyset_string)
		// Create the actor keyset
		if keyset_string == "" {
			log.Fatal("You need to define actorKeyset or generate a new one.")
		}
		unpackedActorKeyset, err := set.Unpack(keyset_string)
		if err != nil {
			log.Fatalf("Failed to unpack keyset: %v", err)
		}
		keyset = unpackedActorKeyset

		log.Debug("Unpacked keyset and set it to actor.")
	}
}

func GetKeyset() *set.Keyset {
	return keyset
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

func GetForcePublish() bool {
	return *forcePublish
}

func GetDiscoveryTimeout() time.Duration {
	return time.Duration(discoveryTimeout) * time.Second
}
