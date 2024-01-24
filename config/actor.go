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
func InitIdentity() {

	keyset_string := viper.GetString("actor.keyset")
	nick := viper.GetString("actor.nick")

	// Generate a new keysets if requested
	if viper.GetBool("generate") {

		log.Debugf("config.initIdentity: Generating new keyset for %s", nick)
		keyset_string = generateKeyset(nick)
		fmt.Println(keyset_string)

		keyset, err := set.Unpack(keyset_string)
		if err != nil {
			log.Errorf("config.initIdentity: Failed to unpack keyset: %v", err)
			os.Exit(65) // EX_DATAERR
		}

		if viper.GetBool("publish") {
			publishIdentity(keyset)
		}

		os.Exit(0)
	}

	log.Debugf("config.initIdentity: %s", keyset_string)
	// Create the actor keyset
	if keyset_string == "" {
		log.Errorf("config.initIdentity: You need to define actorKeyset or generate a new one.")
		os.Exit(64) // EX_USAGE
	}

	keyset, err := set.Unpack(viper.GetString("actor.keyset"))
	if err != nil {
		log.Errorf("config.initIdentity: Failed to unpack keyset: %v", err)
		os.Exit(65) // EX_DATAERR
	}

	if viper.GetBool("publish") {
		if GetKeysetString() != "" {
			publishIdentity(keyset)
		} else {
			log.Errorf("No actor keyset to publish.")
		}
	}

	viper.Set("keyset", keyset)

}

func generateKeyset(nick string) string {

	ks, err := set.GetOrCreate(nick)
	if err != nil {
		log.Errorf("Failed to generate new keyset: %v", err)
		os.Exit(70) // EX_SOFTWARE
	}
	log.Debugf("Created new keyset: %v", ks)

	pks, err := ks.Pack()
	if err != nil {
		log.Errorf("Failed to pack keyset: %v", err)
		os.Exit(70) // EX_SOFTWARE
	}
	log.Debugf("Packed keyset: %v", pks)

	return pks
}

func publishIdentity(k *set.Keyset) {

	d, err := doc.NewFromKeyset(k)
	if err != nil {
		log.Errorf("config.publishIdentity: failed to create DOC: %v", err)
		os.Exit(75) // EX_TEMPFAIL
	}

	assertionMethod, err := d.GetAssertionMethod()
	if err != nil {
		log.Errorf("config.publishIdentity: failed to get verification method: %v", err)
		os.Exit(75) // EX_TEMPFAIL
	}
	d.Sign(k.SigningKey, assertionMethod)

	_, err = d.Publish(nil)
	if err != nil {
		log.Errorf("config.publishIdentity: failed to publish DOC: %v", err)
		os.Exit(75) // EX_TEMPFAIL
	}
	log.Debugf("config.publishIdentity: published DOC: %s", d.ID)

}

func GetKeyset() *set.Keyset {
	return viper.Get("keyset").(*set.Keyset)
}

func GetIPFSKey() *ipfs.Key {
	return GetKeyset().IPFSKey
}

func GetKeysetString() string {
	return viper.GetString("actor.keyset")
}

func GetNick() string {
	return viper.GetString("actor.nick")
}

func GetPublish() bool {

	return viper.GetBool("publish")
}

func GetHome() string {
	return viper.GetString("location.home")
}
