package config

import (
	"fmt"
	"os"

	"github.com/bahner/go-ma/did/doc"
	"github.com/bahner/go-ma/key/set"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func initKeyset(keyset_string string) {

	// Generate a new keysets if requested
	if *generate || *genenv {
		log.Debugf("config.initKeyset: Generating new keyset for %s", nick)
		keyset_string = generateKeyset()
	}

	log.Debugf("config.initKeyset: %s", keyset_string)
	// Create the actor keyset
	if keyset_string == "" {
		log.Fatal("config.initKeyset: You need to define actorKeyset or generate a new one.")
	}

	keyset, err = set.Unpack(keyset_string)
	if err != nil {
		log.Fatalf("config.initKeyset: Failed to unpack keyset: %v", err)
	}

	if *publish || *forcePublish {
		if keyset_string != "" {
			publishIdentity(keyset)
		} else {
			log.Errorf("No actor keyset to publish.")
		}
	}

	if *genenv {
		os.Exit(0)
	}

}

func generateKeyset() string {

	if nick == "ghost" {
		log.Fatal("You need to set a nick when generating an identity.")
	}

	ks, err := set.New(nick, *forcePublish)
	if err != nil {
		log.Fatalf("Failed to generate new keyset: %v", err)
	}

	pks, err := ks.Pack()
	if err != nil {
		log.Fatalf("Failed to pack keyset: %v", err)
	}

	if *genenv {
		fmt.Println("export " + keyset_var + "=" + pks)
	}

	return pks
}

func publishIdentity(k *set.Keyset) {

	err := k.IPNSKey.ExportToIPFS(*forcePublish)
	if err != nil {
		log.Debugf(errors.Cause(err).Error())
		log.Fatalf("config.publishIdentity: failed to export keyset: %v", err)
	}
	log.Infof("Exported IPNSkey to IPFS: %s", k.IPNSKey.DID)

	d, err := doc.NewFromKeyset(keyset, k.IPNSKey.DID)
	if err != nil {
		log.Fatalf("config.publishIdentity: failed to create DOC: %v", err)
	}

	assertionMethod, err := d.GetAssertionMethod()
	if err != nil {
		log.Fatalf("config.publishIdentity: failed to get verification method: %v", err)
	}
	d.Sign(k.SigningKey, assertionMethod)

	_, err = d.Publish()
	if err != nil {
		log.Fatalf("config.publishIdentity: failed to publish DOC: %v", err)
	}
	log.Debugf("config.publishIdentity: published DOC: %s", d.ID)

}

func GetKeyset() *set.Keyset {
	return keyset
}
