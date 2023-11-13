package config

import (
	"flag"
	"os"
	"time"

	"github.com/bahner/go-ma/key/set"
	nanoid "github.com/matoous/go-nanoid/v2"
	log "github.com/sirupsen/logrus"
	"go.deanishe.net/env"
)

const (
	actor_keyset_var        = "GO_HOME_ACTOR_KEYSET"
	room_keyset_var         = "GO_HOME_ROOM_KEYSET"
	defaultDiscoveryTimeout = 300
)

var (
	discoveryTimeout int = env.GetInt("GO_HOME_DISCOVERY_TIMEOUT", defaultDiscoveryTimeout)

	logLevel            string = env.Get("GO_HOME_LOG_LEVEL", "error")
	nick                string = env.Get("USER")
	actor_keyset_string string = env.Get(actor_keyset_var, "")
	room_keyset_string  string = env.Get(room_keyset_var, "")

	generate     *bool
	genenv       *bool
	publish      *bool
	forcePublish *bool

	RoomKeyset  *set.Keyset
	ActorKeyset *set.Keyset

	randomRoomNick, _ = nanoid.New()
)

func Init() {

	// Flags - user configurations
	flag.StringVar(&logLevel, "loglevel", logLevel, "Loglevel to use for application")
	flag.IntVar(&discoveryTimeout, "discoveryTimeout", discoveryTimeout, "Timeout for peer discovery")

	// Actor
	flag.StringVar(&nick, "nick", nick, "Nickname to use in character creation")
	flag.StringVar(&actor_keyset_string, "actorKeyset", actor_keyset_string, "Base58 encoded secret key used to identify the client. You.")

	// Room
	flag.StringVar(&room_keyset_string, "roomKeyset", room_keyset_string, "Base58 encoded secret key used to identify your room.")

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
		actor_keyset_string = generateKeyset(actor_keyset_var, nick, *forcePublish)
		room_keyset_string = generateKeyset(room_keyset_var, randomRoomNick, *forcePublish)
	}

	if *publish || *forcePublish {
		if actor_keyset_string != "" {
			publishKeyset(ActorKeyset, *forcePublish)
		} else {
			log.Errorf("No actor keyset to publish.")
		}

		if room_keyset_string != "" {
			publishKeyset(RoomKeyset, *forcePublish)
		} else {
			log.Errorf("No room keyset to publish.")
		}
	}

	if *genenv || *generate {
		os.Exit(0)
	}

	log.Debugf("actor_keyset_string: %s", actor_keyset_string)
	// Create the actor keyset
	if actor_keyset_string == "" {
		log.Fatal("You need to define actorKeyset or generate a new one.")
	}
	unpackedActorKeyset, err := set.Unpack(actor_keyset_string)
	if err != nil {
		log.Fatalf("Failed to unpack keyset: %v", err)
	}
	ActorKeyset = unpackedActorKeyset

	// Create the room keyset
	if room_keyset_string == "" {
		log.Fatal("You need to define roomKeyset or generate a new one.")
	}
	unpackedRoomKeyset, err := set.Unpack(room_keyset_string)
	if err != nil {
		log.Fatalf("Failed to unpack keyset: %v", err)
	}
	RoomKeyset = unpackedRoomKeyset

	log.Debug("Unpacked keyset and set it to actor.")
}

func GetActorKeyset() *set.Keyset {
	return ActorKeyset
}

func GetRoomKeyset() *set.Keyset {
	return RoomKeyset
}

func GetNick() string {
	return nick
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
