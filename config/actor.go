package config

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/bahner/go-ma/key/set"
	log "github.com/sirupsen/logrus"
)

const (
	defaultLocation   string = "did:ma:k2k4r8p5sxlnznc9ral4fueapazs36hqoj3go1g0o7662gnk9skhfrik#pong"
	fakeActorIdentity string = "NO_DEFAULT_ACTOR_IDENITY"
)

var (
	keyset           set.Keyset
	ErrEmptyIdentity = fmt.Errorf("identity is empty")
	ErrFakeIdentity  = fmt.Errorf("your identity is fake. You need to define actorKeyset or generate a new one")
	ErrEmptyNick     = fmt.Errorf("nick is empty")
)

func InitActorFlags() {

	pflag.StringP("nick", "n", "", "Nickname to use in character creation")
	pflag.StringP("location", "l", defaultLocation, "DID of the location to visit")

	viper.BindPFlag("actor.nick", pflag.Lookup("nick"))
	viper.BindPFlag("actor.location", pflag.Lookup("location"))

	viper.SetDefault("actor.location", defaultLocation)
	viper.SetDefault("actor.nick", defaultNick())

}

// Set the default nick to the user's username, unless a profile is set.
func defaultNick() string {

	if Profile() == defaultProfile {
		return os.Getenv("USER")
	}

	return Profile()
}

// Load a keyset from string and initiate an Actor.
// This is optional, but if you want to use the actor package, you need to call this.
func InitActor() {

	keyset_string := actorIdentity()
	if keyset_string == fakeActorIdentity {
		panic(ErrFakeIdentity)
	}

	log.Debugf("config.initActor: %s", keyset_string)
	// Create the actor keyset
	if keyset_string == "" {
		panic(ErrEmptyIdentity.Error())
	}

	// This function fails fatally, so no return value
	initActorKeyset(keyset_string)

	if PublishFlag() && keyset_string != "" {
		fmt.Println("Publishing identity to IPFS...")
		err := PublishIdentityFromKeyset(keyset)
		if err != nil {
			log.Warnf("config.initActor: %v", err)
		}
	}

}

// Fetches the actor nick from the config or the command line
// NB! This is a little more complex than the other config functions, as it
// needs to fetch the nick from the command line if it's not in the config.
// Due to being a required parameter when generating a new keyset.
func ActorNick() string {

	return viper.GetString("actor.nick")

}

func ActorLocation() string {

	return viper.GetString("actor.location")
}

func ActorKeyset() set.Keyset {
	return keyset
}
func actorIdentity() string {

	return viper.GetString("actor.identity")

}

func initActorKeyset(keyset_string string) {

	log.Debugf("config.initActor: %s", keyset_string)
	// Create the actor keyset
	if keyset_string == "" {
		log.Errorf("config.initActor: You need to define actorKeyset or generate a new one.")
		os.Exit(64) // EX_USAGE
	}

	var err error

	keyset, err = set.Unpack(keyset_string)
	if err != nil {
		log.Errorf("config.initActor: %v", err)
		os.Exit(70) // EX_SOFTWARE
	}

}
