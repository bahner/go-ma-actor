package config

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/bahner/go-ma/did/doc"
	"github.com/bahner/go-ma/key/ipfs"
	"github.com/bahner/go-ma/key/set"
	log "github.com/sirupsen/logrus"
)

const defaultNick string = "ghost"

func init() {
	// Keyset
	pflag.BoolP("generate", "g", false, "Generates a new keyset")
	viper.BindPFlag("generate", pflag.Lookup("generate"))

	pflag.BoolP("publish", "p", false, "Publishes keyset to IPFS")
	viper.BindPFlag("publish", pflag.Lookup("publish"))

	// Nick used for keyset generation (fragment)
	pflag.StringP("nick", "n", defaultNick, "Nickname to use in character creation")
	viper.BindPFlag("actor.nick", pflag.Lookup("nick"))
	err := viper.BindEnv("actor.nick", "USER")
	if err != nil {
		log.Fatalf("Error binding environment variable 'USER': %s\n", err)
	}

	pflag.StringP("location", "l", "", "DID of the initial location.")
	viper.BindPFlag("location.home", pflag.Lookup("home"))

}

// Load a keyset from string and initiate an Actor.
// This is optional, but if you want to use the actor package, you need to call this.
func InitActor() {

	keyset_string := viper.GetString("actor.identity")

	log.Debugf("config.initIdentity: %s", keyset_string)
	// Create the actor keyset
	if keyset_string == "" {
		log.Errorf("config.initIdentity: You need to define actorKeyset or generate a new one.")
		os.Exit(64) // EX_USAGE
	}

	keyset, err := set.Unpack(keyset_string)
	if err != nil {
		log.Errorf("config.initIdentity: Failed to unpack keyset: %v", err)
	}

	if viper.GetBool("publish") && keyset_string != "" {
		err := publishIdentityFromKeyset(keyset)
		if err != nil {
			log.Errorf("config.initIdentity: Failed to publish keyset: %v", err)
			os.Exit(75) // EX_TEMPFAIL
		}
	}

	viper.Set("keyset", keyset)

}

func handleGenerateOrExit() {
	// Generate a new keysets if requested

	keyset_string, err := generateAndPrintActorIdentity()
	if err != nil {
		log.Errorf("config.initIdentity: Failed to generate keyset: %v", err)
		os.Exit(70) // EX_SOFTWARE
	}

	if viper.GetBool("publish") {
		err = publishActorIdentityFromString(keyset_string)
		if err != nil {
			log.Errorf("config.initIdentity: Failed to publish keyset: %v", err)
			os.Exit(75) // EX_TEMPFAIL
		}
	}

	err = generateAndPrintNodeIdentity()
	if err != nil {
		log.Errorf("config.initIdentity: Failed to generate node identity: %v", err)
		os.Exit(70) // EX_SOFTWARE
	}

}

func generateAndPrintActorIdentity() (string, error) {

	nick := viper.GetString("actor.nick")

	keyset_string, err := generateKeyset(nick)
	if err != nil {
		return "", fmt.Errorf("config.initIdentity: Failed to generate keyset: %v", err)
	}

	fmt.Println(ENV_PREFIX + "_ACTOR_IDENTITY=" + keyset_string)

	return keyset_string, nil
}

func publishActorIdentityFromString(keyset_string string) error {

	keyset, err := set.Unpack(keyset_string)
	if err != nil {
		log.Errorf("config.initIdentity: Failed to unpack keyset: %v", err)
	}

	err = publishIdentityFromKeyset(keyset)
	if err != nil {
		return fmt.Errorf("config.initIdentity: Failed to publish keyset: %v", err)
	}

	return nil
}

func generateKeyset(nick string) (string, error) {

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

func publishIdentityFromKeyset(k *set.Keyset) error {

	d, err := doc.NewFromKeyset(k)
	if err != nil {
		return fmt.Errorf("config.publishIdentity: failed to create DOC: %v", err)
	}

	assertionMethod, err := d.GetAssertionMethod()
	if err != nil {
		return fmt.Errorf("config.publishIdentity: failed to get verification method: %v", err)
	}
	d.Sign(k.SigningKey, assertionMethod)

	_, err = d.Publish(nil)
	if err != nil {
		return fmt.Errorf("config.publishIdentity: failed to publish DOC: %v", err)

	}
	log.Debugf("config.publishIdentity: published DOC: %s", d.ID)

	return nil
}

func GetKeyset() *set.Keyset {
	return viper.Get("keyset").(*set.Keyset)
}

func GetIPFSKey() *ipfs.Key {
	return GetKeyset().IPFSKey
}

func GetActorIdentity() string {
	return viper.GetString("actor.identity")
}

func GetActorNick() string {
	return viper.GetString("actor.nick")
}

func GetPublish() bool {

	return viper.GetBool("publish")
}

func GetHome() string {
	return viper.GetString("location.home")
}
