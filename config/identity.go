package config

import (
	"flag"
	"os"

	"github.com/bahner/go-ma/did/doc"
	"github.com/bahner/go-ma/key/ipfs"
	"github.com/bahner/go-ma/key/set"
	log "github.com/sirupsen/logrus"
	"go.deanishe.net/env"
)

var (

	// Booleans with control flow
	generate = flag.Bool("generate", false, "Generates a new keyset")
	publish  = flag.Bool("publish", false, "Publishes keyset to IPFS when using genenv or generate")

	// Entities
	nick          = flag.String("nick", env.Get("USER", defaultNick), "Nickname to use in character creation")
	keyset        *set.Keyset
	keyset_string = flag.String("keyset", env.Get(GO_MA_ACTOR_IDENTITY_VAR, defaultKeyset),
		"Base58 encoded *secret* keyset used to identify the client. You. You can use environment variable "+GO_MA_ACTOR_IDENTITY_VAR+" to set this.")
)

func InitIdentity() {

	var err error

	// Generate a new keysets if requested
	if *generate {
		log.Debugf("config.initIdentity: Generating new keyset for %s", *nick)
		*keyset_string = generateKeyset()
	}

	log.Debugf("config.initIdentity: %s", *keyset_string)
	// Create the actor keyset
	if *keyset_string == "" {
		log.Errorf("config.initIdentity: You need to define actorKeyset or generate a new one.")
		os.Exit(64) // EX_USAGE
	}

	keyset, err = set.Unpack(*keyset_string)
	if err != nil {
		log.Errorf("config.initIdentity: Failed to unpack keyset: %v", err)
		os.Exit(65) // EX_DATAERR
	}

	if *publish {
		if *keyset_string != "" {
			publishIdentity(keyset)
		} else {
			log.Errorf("No actor keyset to publish.")
		}
	}

}

func generateKeyset() string {

	if *nick == "ghost" {
		log.Errorf("You need to set a nick when generating an identity.")
		os.Exit(64) // EX_USAGE
	}

	// ks, err := set.New(*nick, *forcePublish)
	ks, err := set.GetOrCreate(*nick)
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

	d, err := doc.NewFromKeyset(keyset)
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
	return keyset
}

func GetIPFSKey() *ipfs.Key {
	return keyset.IPFSKey
}

func GetIdentityString() string {
	return *keyset_string
}

func GetNick() string {
	return *nick
}

func GetPublish() bool {

	return *publish
}
