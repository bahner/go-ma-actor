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
	configActorFlagset    = pflag.NewFlagSet("actor", pflag.ExitOnError)
	actorKeyset           set.Keyset
	configActorFlagsOnce  sync.Once
	configActorKeysetLoad sync.Once
	actorNick             string
	actorLocation         string
	actorMnemonic         string
	actorKeysetPath       string
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
	Keyset   string `yaml:"keyset-path"`
	Location string `yaml:"location"`
	Mnemonic string `yaml:"mnemonic"`
	Nick     string `yaml:"nick"`
}

// Config for actor. Remember to parse the flags first.
// Eg. ActorFlags()
func Actor() ActorConfig {

	// If we are generating a new identity we should publish it
	if GenerateFlag() {
		publishDIDDocumentFromKeyset(actorKeyset)
	}

	return ActorConfig{
		Keyset:   getActorKeysetPath(),
		Location: ActorLocation(),
		Mnemonic: ActorMnemonic(),
		Nick:     ActorNick(),
	}
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

func ActorKeyset() set.Keyset {
	configActorKeysetLoad.Do(func() {

		ctx, cancel := context.WithTimeout(context.Background(), configActorIPFSTimeout)
		defer cancel()

		actorKeyset, err = set.LoadFromIPFS(ctx, getActorKeysetPath(), ActorMnemonic())
		if err != nil {
			log.Fatalf("config.Actor: %v", err)
		}
	})

	return actorKeyset
}

func getActorKeysetPath() string {

	if GenerateFlag() && actorKeysetPath == "" {
		actorKeysetPath, err = generateKeysetPath(ActorNick())
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

// Generates a new keyset and returns the keyset as a string
func generateKeysetPath(nick string) (string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), configActorIPFSTimeout)
	defer cancel()

	privKey, err := getOrCreateIdentity(nick)
	if err != nil {
		return "", fmt.Errorf("failed to generate identity: %w", err)
	}
	keyset, err := set.New(privKey, nick)
	if err != nil {
		return "", fmt.Errorf("failed to generate new keyset: %w", err)
	}
	log.Debugf("Created new keyset: %v", keyset)

	cid, err := keyset.SaveToIPFS(ctx, ActorMnemonic())
	if err != nil {
		return "", fmt.Errorf("failed to save keyset to IPFS: %w", err)
	}

	log.Debugf("Keyset saved to IPFS: %s\n", cid.String())
	return cid.String(), nil
}

func publishDIDDocumentFromKeyset(k set.Keyset) error {

	d, err := doc.NewFromKeyset(k)
	if err != nil {
		return fmt.Errorf("config.publishDIDDocumentFromKeyset: failed to create DOC: %w", err)
	}

	_, err = d.Publish()
	if err != nil {
		return fmt.Errorf("config.publishDIDDocumentFromKeyset: %w", err)

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
