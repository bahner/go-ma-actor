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

	var err error

	// Generate a new keysets if requested
	if *generate || *genenv {
		log.Debugf("config.initKeyset: Generating new keyset for %s", nick)
		keyset_string = generateKeyset()
	}

	log.Debugf("config.initKeyset: %s", keyset_string)
	// Create the actor keyset
	if keyset_string == "" {
		log.Errorf("config.initKeyset: You need to define actorKeyset or generate a new one.")
		os.Exit(64) // EX_USAGE
	}

	keyset, err = set.Unpack(keyset_string)
	if err != nil {
		log.Errorf("config.initKeyset: Failed to unpack keyset: %v", err)
		os.Exit(65) // EX_DATAERR
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
		log.Errorf("You need to set a nick when generating an identity.")
		os.Exit(64) // EX_USAGE
	}

	ks, err := set.New(nick, *forcePublish)
	if err != nil {
		log.Errorf("Failed to generate new keyset: %v", err)
		os.Exit(70) // EX_SOFTWARE
	}

	pks, err := ks.Pack()
	if err != nil {
		log.Errorf("Failed to pack keyset: %v", err)
		os.Exit(70) // EX_SOFTWARE
	}

	if *genenv {
		fmt.Println("export " + GO_MA_ACTOR_KEYSET_VAR + "=" + pks)
	}

	return pks
}

func publishIdentity(k *set.Keyset) {

	err := k.IPNSKey.ExportToIPFS(*forcePublish)
	if err != nil {
		log.Debugf(errors.Cause(err).Error())
		log.Errorf("config.publishIdentity: failed to export keyset: %v", err)
		os.Exit(75) // EX_TEMPFAIL
	}
	log.Infof("Exported IPNSkey to IPFS: %s", k.IPNSKey.DID)

	d, err := doc.NewFromKeyset(keyset, k.IPNSKey.DID)
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

	_, err = d.Publish()
	if err != nil {
		log.Errorf("config.publishIdentity: failed to publish DOC: %v", err)
		os.Exit(75) // EX_TEMPFAIL
	}
	log.Debugf("config.publishIdentity: published DOC: %s", d.ID)

}

func GetKeyset() *set.Keyset {
	return keyset
}
