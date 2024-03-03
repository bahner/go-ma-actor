package config

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/bahner/go-ma/did/doc"
	"github.com/bahner/go-ma/key/set"
	log "github.com/sirupsen/logrus"
)

const (
	defaultActor string = "actor"
)

var keyset set.Keyset

func init() {
	// Keyset
	pflag.BoolP("generate", "g", false, "Generates a new keyset")

	pflag.BoolP("publish", "p", false, "Publishes keyset to IPFS")

	// Nick used for keyset generation (fragment)
	pflag.StringP("nick", "n", defaultActor, "Nickname to use in character creation")
	viper.BindPFlag("actor.nick", pflag.Lookup("nick"))

	pflag.StringP("location", "l", "", "DID of the initial location.")
	viper.BindPFlag("actor.home", pflag.Lookup("home"))

}

// Load a keyset from string and initiate an Actor.
// This is optional, but if you want to use the actor package, you need to call this.
func InitActor() {

	keyset_string := viper.GetString("actor.identity")

	log.Debugf("config.initActor: %s", keyset_string)
	// Create the actor keyset
	if keyset_string == "" {
		log.Errorf("config.initActor: You need to define actorKeyset or generate a new one.")
		os.Exit(64) // EX_USAGE
	}

	// This function fails fatally, so no return value
	initActorKeyset()

	if publishFlag() && keyset_string != "" {
		fmt.Print("Publishing identity to IPFS...")
		err := publishIdentityFromKeyset(keyset)
		if err != nil {
			log.Warnf("config.initActor: %v", err)
		}
		fmt.Println("done.")
	}

}

func ActorNick() string {
	return viper.GetString("actor.nick")
}

func ActorHome() string {
	return viper.GetString("actor.home")
}

func ActorDid() string {
	return viper.GetString("actor.did")
}

func ActorKeyset() set.Keyset {
	return keyset
}

// Genreates a libp2p and actor identity and returns the keyset and the actor identity
// These are imperative, so failure to generate them is a fatal error.
func handleGenerateOrExit() (string, string) {

	// Generate a new keysets if requested
	nick := viper.GetString("actor.nick")

	keyset_string, err := generateKeysetString(nick)
	if err != nil {
		log.Errorf("config.handleGenerateOrExit: %v", err)
		os.Exit(70) // EX_SOFTWARE
	}

	ni, err := generateNodeIdentity()
	if err != nil {
		log.Errorf("config.handleGenerateOrExit: %v", err)
		os.Exit(70) // EX_SOFTWARE
	}

	if publishFlag() {
		err = publishActorIdentityFromString(keyset_string)
		if err != nil {
			log.Warnf("config.handleGenerateOrExit: %v", err)
		}
	}

	return keyset_string, ni
}

func publishActorIdentityFromString(keyset_string string) error {

	keyset, err := set.Unpack(keyset_string)
	if err != nil {
		log.Errorf("config.publishActorIdentityFromString: Failed to unpack keyset: %v", err)
	}

	err = publishIdentityFromKeyset(keyset)
	if err != nil {
		return fmt.Errorf("config.publishActorIdentityFromString: Failed to publish keyset: %v", err)
	}

	return nil
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

func publishIdentityFromKeyset(k set.Keyset) error {

	d, err := doc.NewFromKeyset(k)
	if err != nil {
		return fmt.Errorf("config.publishIdentityFromKeyset: failed to create DOC: %v", err)
	}

	assertionMethod, err := d.GetAssertionMethod()
	if err != nil {
		return fmt.Errorf("config.publishIdentityFromKeyset: %w", err)
	}
	d.Sign(k.SigningKey, assertionMethod)

	// Publication options
	opts := docPublishOptions()

	_, err = d.Publish(opts)
	if err != nil {
		return fmt.Errorf("config.publishIdentityFromKeyset: %w", err)

	}
	log.Debugf("Published identity: %s", d.ID)

	return nil
}

func initActorKeyset() {

	keyset_string := viper.GetString("actor.identity")

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

func docPublishOptions() *doc.PublishOptions {
	return &doc.PublishOptions{
		Ctx:   context.Background(),
		Force: forceFlag(),
	}
}
