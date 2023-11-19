package config

import (
	"flag"

	"github.com/bahner/go-ma"
	"github.com/bahner/go-ma/key/set"
	nanoid "github.com/matoous/go-nanoid/v2"
	log "github.com/sirupsen/logrus"
	"go.deanishe.net/env"
)

const (
	actor_keyset_var = "GO_HOME_ACTOR_KEYSET"
	room_keyset_var  = "GO_HOME_ROOM_KEYSET"
)

var (
	randomRoomNick, _ = nanoid.New()

	logLevel            string = env.Get("GO_HOME_LOG_LEVEL", "error")
	rendezvous          string = env.Get("GO_HOME_RENDEZVOUS", ma.RENDEZVOUS)
	serviceName         string = env.Get("GO_HOME_SERVICE_NAME", ma.RENDEZVOUS)
	actorNick           string = env.Get("USER")
	actor_keyset_string string = env.Get(actor_keyset_var, "")
	room_keyset_string  string = env.Get(room_keyset_var, "")
	roomNick            string = env.Get("GO_HOME_ROOM_NICK", randomRoomNick)

	generate     *bool
	genenv       *bool
	publish      *bool
	forcePublish *bool

	RoomKeyset  *set.Keyset
	ActorKeyset *set.Keyset
)

func Init() {

	// Flags - user configurations
	flag.StringVar(&logLevel, "loglevel", logLevel, "Loglevel to use for application")
	flag.StringVar(&rendezvous, "rendezvous", rendezvous, "Unique string to identify group of nodes. Share this with your friends to let them connect with you")
	flag.StringVar(&serviceName, "servicename", serviceName, "serviceName to use for MDNS discovery")

	// Actor
	flag.StringVar(&actorNick, "actorNick", actorNick, "Nickname to use in character creation")
	flag.StringVar(&actor_keyset_string, "actor_keyset", actor_keyset_string, "Base58 encoded secret key used to identify the client. You.")

	// Room
	flag.StringVar(&roomNick, "roomNick", roomNick, "Nickname to use in room creation")
	flag.StringVar(&room_keyset_string, "room_keyset", actor_keyset_string, "Base58 encoded secret key used to identify your room.")

	// Booleans with control flow
	generate = flag.Bool("generate", false, "Generates one-time keyset and uses it")
	genenv = flag.Bool("genenv", false, "Generates a keyset and prints it to stdout and uses it")
	publish = flag.Bool("publish", false, "Publishes keyset to IPFS when using genenv or generate")
	forcePublish = flag.Bool("force-publish", false, "Force publish even if keyset is already published")

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
		actor_keyset_string = generateKeyset(actor_keyset_var, actorNick)
		room_keyset_string = generateKeyset(room_keyset_var, randomRoomNick)
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
	ActorKeyset = &unpackedActorKeyset

	// Create the room keyset
	if room_keyset_string == "" {
		log.Fatal("You need to define roomKeyset or generate a new one.")
	}
	unpackedRoomKeyset, err := set.Unpack(room_keyset_string)
	if err != nil {
		log.Fatalf("Failed to unpack keyset: %v", err)
	}
	RoomKeyset = &unpackedRoomKeyset

	// Publish the keysets if requested
	if *publish || *forcePublish {
		publishKeyset(ActorKeyset)
		publishKeyset(RoomKeyset)
	}

	log.Debug("Unpacked keyset and set it to actor.")
}

func GetActorKeyset() *set.Keyset {
	return ActorKeyset
}

func GetRoomKeyset() *set.Keyset {
	return RoomKeyset
}

func GetRendezvous() string {
	return rendezvous
}

func GetServiceName() string {
	return serviceName
}

func GetActorNick() string {
	return actorNick
}

func GetRoomNick() string {
	return roomNick
}

func GetLogLevel() string {
	return logLevel
}

func GetForcePublish() bool {
	return *forcePublish
}
