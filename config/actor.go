package config

// This file contains the configuration for the actor package.
// It also somewhat strangeky initialises the identoty and generates a new one if needed.
// This is because it's so low level and the identity is needed for the keyset.

import (
	"fmt"
	"os"
	"sync"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/bahner/go-ma/did/doc"
	"github.com/bahner/go-ma/key/set"
	log "github.com/sirupsen/logrus"
)

const (
	defaultLocation   string = "did:ma:k2k4r8p5sxlnznc9ral4fueapazs36hqoj3go1g0o7662gnk9skhfrik#pong"
	fakeActorIdentity string = "NO_DEFAULT_ACTOR_IDENITY"
)

var (
	actorFlagset     = pflag.NewFlagSet("actor", pflag.ExitOnError)
	actorKeyset      set.Keyset
	actorFlagsOnce   sync.Once
	ErrEmptyIdentity = fmt.Errorf("identity is empty")
	ErrEmptyNick     = fmt.Errorf("nick is empty")
	ErrFakeIdentity  = fmt.Errorf("your identity is fake. You need to define actorKeyset or generate a new one")
	nick             string
	location         string
)

// Initialise command line flags for the actor package
// The actor is optional for some commands, but required for others.
// exitOnHelp means that this function is the last called when help is needed.
// and the program should exit.
func actorFlags() {

	actorFlagsOnce.Do(func() {

		actorFlagset.StringVarP(&nick, "nick", "n", "", "Nickname to use in character creation")
		actorFlagset.StringVarP(&location, "location", "l", defaultLocation, "DID of the location to visit")

		viper.BindPFlag("actor.nick", actorFlagset.Lookup("nick"))
		viper.BindPFlag("actor.location", actorFlagset.Lookup("location"))

		viper.SetDefault("actor.location", defaultLocation)
		viper.SetDefault("actor.nick", defaultNick())

		if HelpNeeded() {
			fmt.Println("Actor Flags:")
			actorFlagset.PrintDefaults()

		}
	})
}

type ActorConfig struct {
	Identity string `yaml:"identity"`
	Nick     string `yaml:"nick"`
	Location string `yaml:"location"`
}

// Config for actor. Remember to parse the flags first.
// Eg. ActorFlags()
func Actor() ActorConfig {

	// Fetch the identity from the config or generate one
	identity, err := actorIdentity()
	if err != nil {
		panic(err)
	}

	// Unpack the keyset from the identity
	initActorKeyset(identity)

	// If we are generating a new identity we should publish it
	if GenerateFlag() {
		publishIdentityFromKeyset(actorKeyset)
	}

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

	// This is used early, so command line takes precedence
	if actorFlagset.Lookup("nick").Changed {
		return actorFlagset.Lookup("nick").Value.String()
	}
	return viper.GetString("actor.nick")
}

func ActorLocation() string {
	return viper.GetString("actor.location")
}

func ActorKeyset() set.Keyset {
	return actorKeyset
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

	actorKeyset, err = set.Unpack(keyset_string)
	if err != nil {
		log.Errorf("config.initActor: %v", err)
		os.Exit(70) // EX_SOFTWARE
	}
}

// Generates a new keyset and returns the keyset as a string
func generateKeysetString(nick string) (string, error) {

	privKey, err := getOrCreateIdentity(nick)
	if err != nil {
		return "", fmt.Errorf("failed to generate identity: %w", err)
	}
	keyset, err := set.New(privKey, nick)
	if err != nil {
		return "", fmt.Errorf("failed to generate new keyset: %w", err)
	}
	log.Debugf("Created new keyset: %v", keyset)

	pks, err := keyset.Pack()
	if err != nil {
		return "", fmt.Errorf("failed to pack keyset: %w", err)
	}
	log.Debugf("Packed keyset: %v", pks)

	return pks, nil
}

func publishIdentityFromKeyset(k set.Keyset) error {

	d, err := doc.NewFromKeyset(k)
	if err != nil {
		return fmt.Errorf("config.publishIdentityFromKeyset: failed to create DOC: %w", err)
	}

	assertionMethod, err := d.GetAssertionMethod()
	if err != nil {
		return fmt.Errorf("config.publishIdentityFromKeyset: %w", err)
	}
	d.Sign(k.SigningKey, assertionMethod)

	_, err = d.Publish()
	if err != nil {
		return fmt.Errorf("config.publishIdentityFromKeyset: %w", err)

	}
	log.Debugf("Published identity: %s", d.ID)

	return nil
}

func getOrCreateIdentity(name string) (crypto.PrivKey, error) {

	hasKey, err := Keystore().Has(name)
	if hasKey {
		return Keystore().Get(name)
	}
	if err != nil {
		log.Errorf("failed to get private key from keystore: %s", err)
		return nil, err
	}

	privKey, _, err := crypto.GenerateKeyPair(crypto.Ed25519, -1)
	if err != nil {
		log.Errorf("failed to generate node identity: %s", err)
		return nil, err
	}

	err = Keystore().Put(name, privKey)
	if err != nil {
		log.Errorf("failed to store private key to keystore: %s", err)
		return nil, err
	}

	return privKey, nil

}
