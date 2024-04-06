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

// Initialise command line flags for the actor package
// The actor is optional for some commands, but required for others.
func ActorFlags() {

	pflag.StringP("nick", "n", "", "Nickname to use in character creation")
	pflag.StringP("location", "l", defaultLocation, "DID of the location to visit")

	viper.BindPFlag("actor.nick", pflag.Lookup("nick"))
	viper.BindPFlag("actor.location", pflag.Lookup("location"))

	viper.SetDefault("actor.location", defaultLocation)
	viper.SetDefault("actor.nick", defaultNick())

}

type ActorConfig struct {
	Identity string `yaml:"identity"`
	Nick     string `yaml:"nick"`
	Location string `yaml:"location"`
}

// Cofig for actor. Remember to parse the flags first.
// Eg. ActorFlags()
func Actor() ActorConfig {

	identity, err := actorIdentity()
	if err != nil {
		panic(err)
	}

	initActorKeyset(identity)

	return ActorConfig{
		Identity: identity,
		Nick:     ActorNick(),
		Location: ActorLocation(),
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

func actorIdentity() (string, error) {

	if GenerateFlag() {
		return generateKeysetString(ActorNick())
	}

	return viper.GetString("actor.identity"), nil
}

// Set the default nick to the user's username, unless a profile is set.
func defaultNick() string {

	if Profile() == defaultProfile {
		return os.Getenv("USER")
	}

	return Profile()
}

func initActorKeyset(keyset_string string) {

	var err error

	log.Debugf("config.initActor: %s", keyset_string)
	// Create the actor keyset
	if keyset_string == "" {
		log.Errorf("config.initActor: You need to define actorKeyset or generate a new one.")
		os.Exit(64) // EX_USAGE
	}

	keyset, err = set.Unpack(keyset_string)
	if err != nil {
		log.Errorf("config.initActor: %v", err)
		os.Exit(70) // EX_SOFTWARE
	}
}

// Generates a new keyset and returns the keyset as a string
func generateKeysetString(nick string) (string, error) {

	ks, err := set.GetOrCreate(nick)
	if err != nil {
		return "", fmt.Errorf("failed to generate new keyset: %w", err)
	}
	log.Debugf("Created new keyset: %v", ks)

	pks, err := ks.Pack()
	if err != nil {
		return "", fmt.Errorf("failed to pack keyset: %w", err)
	}
	log.Debugf("Packed keyset: %v", pks)

	return pks, nil
}
