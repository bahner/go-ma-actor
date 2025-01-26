package config

// This file contains the configuration for the actor package.
// It also somewhat strangeky initialises the identoty and generates a new one if needed.
// This is because it's so low level and the identity is needed for the keyset.

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/tyler-smith/go-bip39"

	"github.com/bahner/go-ma/did/doc"
	"github.com/bahner/go-ma/key/set"
	log "github.com/sirupsen/logrus"
)

const (
	defaultLocation        string = "did:ma:k2k4r8p5sxlnznc9ral4fueapazs36hqoj3go1g0o7662gnk9skhfrik#pong"
	fakeActorIdentity      string = "NO_DEFAULT_ACTOR_IDENITY"
	configActorIPFSTimeout        = 10 * time.Second
)

var (
	actorDID              string
	actorKeyset           set.Keyset
	actorKeysetPath       string
	actorLocation         string
	actorMnemonic         string
	actorNick             string
	configActorFlagset    = pflag.NewFlagSet("actor", pflag.ExitOnError)
	configActorFlagsOnce  sync.Once
	configActorKeysetLoad sync.Once
)

// Initialise command line flags for the actor package
// The actor is optional for some commands, but required for others.
// exitOnHelp means that this function is the last called when help is needed.
// and the program should exit.
func actorFlags() {

	configActorFlagsOnce.Do(func() {

		configActorFlagset.StringVarP(&actorNick, "nick", "n", "", "Nickname to use in character creation")
		configActorFlagset.StringVarP(&actorLocation, "location", "l", defaultLocation, "DID of the location to visit")
		configActorFlagset.StringVarP(&actorKeysetPath, "keyset-path", "k", "", "IPFS path for keyset to use for the actor")
		configActorFlagset.StringVarP(&actorMnemonic, "mnemonic", "m", "", "BIP-39 Mnemonic to use keyset encryption")

		viper.BindPFlag("actor.nick", configActorFlagset.Lookup("nick"))
		viper.BindPFlag("actor.location", configActorFlagset.Lookup("location"))
		viper.BindPFlag("actor.keyset-path", configActorFlagset.Lookup("keyset-path"))
		viper.BindPFlag("actor.mnemonic", configActorFlagset.Lookup("mnemonic"))

		viper.SetDefault("actor.nick", defaultNick())
		viper.SetDefault("actor.location", defaultLocation)

		if HelpNeeded() {
			fmt.Println("Actor Flags:")
			configActorFlagset.PrintDefaults()

		}
	})
}

type ActorConfig struct {
	DID      string `yaml:"did"`
	Keyset   string `yaml:"keyset-path"`
	Location string `yaml:"location"`
	Mnemonic string `yaml:"mnemonic"`
	Nick     string `yaml:"nick"`
}

// Config for actor. Remember to parse the flags first.
// Eg. ActorFlags()
func Actor() ActorConfig {

	var d *doc.Document

	// If we are generating a new identity we should publish it
	if GenerateFlag() {

		actorKeyset, actorKeysetPath, err = generateActorKeyset(ActorNick())
		if err != nil {
			log.Fatalf("config.ActorKeysetPath: %v", err)
		}

		d, err = publishDIDDocumentFromKeyset(ActorKeyset())
		if err != nil {
			panic(err)
		}
		actorDID = d.ID

	}

	return ActorConfig{
		DID:      ActorDID(),
		Keyset:   ActorKeysetPath(),
		Location: ActorLocation(),
		Mnemonic: ActorMnemonic(),
		Nick:     ActorNick(),
	}
}

func ActorDID() string {
	return viper.GetString("actor.did")
}

// Fetches the actor nick from the config or the command line
// NB! This is a little more complex than the other config functions, as it
// needs to fetch the nick from the command line if it's not in the config.
// Due to being a required parameter when generating a new keyset.
func ActorNick() string {

	// This is used early, so command line takes precedence
	if configActorFlagset.Lookup("nick").Changed {
		return configActorFlagset.Lookup("nick").Value.String()
	}
	return viper.GetString("actor.nick")
}

func ActorLocation() string {
	return viper.GetString("actor.location")
}

func ActorMnemonic() string {
	if GenerateFlag() && actorMnemonic == "" {
		return generateAndSetMnemonic()
	}

	return viper.GetString("actor.mnemonic")
}

// Returns the keyset for the actor after initialisation.
func ActorKeyset() set.Keyset {
	return actorKeyset
}

// Generates and sets the actorKeyset and returns the path to the keyset.
func ActorKeysetPath() string {

	// Not sure if this is too much sugar.
	if GenerateFlag() && actorKeysetPath == "" {
		actorKeyset, actorKeysetPath, err = generateActorKeyset(ActorNick())
		if err != nil {
			log.Fatalf("config.ActorKeysetPath: %v", err)
		}
	}

	return viper.GetString("actor.keyset-path")
}

// Set the default nick to the user's username, unless a profile is set.
func defaultNick() string {

	if Profile() == defaultProfile {
		return os.Getenv("USER")
	}

	return Profile()
}

// Generates and initialises a new keyset for the actor.
// Returns the keyset and the IPFS path to the keyset.
// Also sets the global variables for the helper functions to use.
func generateActorKeyset(nick string) (set.Keyset, string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), configActorIPFSTimeout)
	defer cancel()

	privKey, err := getOrCreateIdentity(nick)
	if err != nil {
		return set.Keyset{}, "", fmt.Errorf("failed to generate identity: %w", err)
	}
	actorKeyset, err = set.New(privKey, nick)
	if err != nil {
		return set.Keyset{}, "", fmt.Errorf("failed to generate new keyset: %w", err)
	}
	log.Debugf("Created new keyset: %v", actorKeyset)

	a, err := actor.New(actorKeyset)
	if err != nil {
		return set.Keyset{}, "", fmt.Errorf("failed to create new actor: %w", err)
	}

	// cid, err := actorKeyset.SaveToIPFS(ctx, ActorMnemonic())
	cid, err := actor.SaveToIPFS(ctx, a)
	if err != nil {
		return set.Keyset{}, "", fmt.Errorf("failed to save keyset to IPFS: %w", err)
	}

	// Set the global _actorKeysetPath
	actorKeysetPath = cid.String()

	log.Debugf("Keyset saved to IPFS: %s\n", cid.String())
	return actorKeyset, actorKeysetPath, nil
}

// Publishes the DID document to the IPFS network and returns the document.
func publishDIDDocumentFromKeyset(k set.Keyset) (*doc.Document, error) {

	d, err := doc.NewFromKeyset(k)
	if err != nil {
		return nil, fmt.Errorf("config.publishDIDDocumentFromKeyset: failed to create DOC: %w", err)
	}

	d.Identity = actorKeysetPath
	err = d.SetTopic(d.ID, doc.DEFAULT_TOPIC_TYPE)
	if err != nil {
		return nil, fmt.Errorf("config.publishDIDDocumentFromKeyset: %w", err)
	}
	err = d.SetP2PHostFromPrivateKey(k.Identity, doc.DEFAULT_HOST_TYPE)
	if err != nil {
		return nil, fmt.Errorf("config.publishDIDDocumentFromKeyset: %w", err)
	}

	_, err = d.Publish()
	if err != nil {
		return nil, fmt.Errorf("config.publishDIDDocumentFromKeyset: %w", err)

	}
	log.Debugf("Published identity: %s", d.ID)

	return d, nil
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

func generateAndSetMnemonic() string {
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		errStr := fmt.Errorf("config.generateMnemonic: %w", err)
		panic(errStr)
	}

	actorMnemonic, err = bip39.NewMnemonic(entropy)
	if err != nil {
		errStr := fmt.Errorf("config.generateMnemonic: %w", err)
		panic(errStr)
	}

	fmt.Printf("Generated new mnemonic: %s\n", actorMnemonic)

	return actorMnemonic
}
