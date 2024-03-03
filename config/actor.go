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

	keyset := ActorKeyset()

	publishFlag, err := pflag.CommandLine.GetBool("publish")
	if err != nil {
		log.Warnf("config.initActor: %v", err)
	}
	if publishFlag && keyset_string != "" {
		fmt.Print("Publishing identity to IPFS...")
		err := publishIdentityFromKeyset(keyset)
		if err != nil {
			log.Warnf("config.initActor: %v", err)
		}
		fmt.Println("done.")
	}

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

	publishFlag, err := pflag.CommandLine.GetBool("publish")
	if err != nil {
		log.Warnf("config.handleGenerateOrExit: %v", err)
	}
	if publishFlag {
		err = publishActorIdentityFromString(keyset_string)
		if err != nil {
			log.Errorf("config.handleGenerateOrExit: %v", err)
			os.Exit(75) // EX_TEMPFAIL
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
		return fmt.Errorf("config.publishIdentityFromKeyset: failed to get verification method: %v", err)
	}
	d.Sign(k.SigningKey, assertionMethod)

	// Publication options
	opts := doc.DefaultPublishOptions()
	opts.Force = viper.GetBool("publish")

	_, err = d.Publish(opts)
	if err != nil {
		return fmt.Errorf("config.publishIdentityFromKeyset: failed to publish DOC: %v", err)

	}
	log.Debugf("config.publishIdentityFromKeyset: published DOC: %s", d.ID)

	return nil
}

func ActorNick() string {
	return viper.GetString("actor.nick")
}

func ActorHome() string {
	return viper.GetString("actor.home")
}

func GetDocPublishOptions() *doc.PublishOptions {
	return &doc.PublishOptions{
		Ctx:   GetPublishContext(),
		Force: viper.GetBool("publish"),
	}
}

func GetPublishContext() context.Context {
	return context.Background()
}

func ActorDid() string {
	return viper.GetString("actor.did")
}

func ActorKeyset() set.Keyset {

	keyset_string := viper.GetString("actor.identity")

	log.Debugf("config.initActor: %s", keyset_string)
	// Create the actor keyset
	if keyset_string == "" {
		log.Errorf("config.initActor: You need to define actorKeyset or generate a new one.")
		os.Exit(64) // EX_USAGE
	}

	keyset, _ := set.Unpack(keyset_string)

	return keyset

}
