package config

import (
	"fmt"
	"os"

	"github.com/bahner/go-ma/did/doc"
	"github.com/bahner/go-ma/key/set"
	log "github.com/sirupsen/logrus"
)

func initKeyset(keyset_string string) {

	// Generate a new keysets if requested
	if *generate || *genenv {
		log.Debugf("generate_keyset: Generating new keyset for %s", nick)
		keyset_string = generateKeyset()
	}

	log.Debugf("actor_keyset_string: %s", keyset_string)
	// Create the actor keyset
	if keyset_string == "" {
		log.Fatal("You need to define actorKeyset or generate a new one.")
	}

	keyset, err = set.Unpack(keyset_string)
	if err != nil {
		log.Fatalf("Failed to unpack keyset: %v", err)
	}

	if *publish || *forcePublish {
		if keyset_string != "" {
			publishKeyset(keyset)
		} else {
			log.Errorf("No actor keyset to publish.")
		}

		if *genenv || *generate {
			os.Exit(0)
		}
		log.Debug("Unpacked keyset and set it to actor.")
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

func publishKeyset(k *set.Keyset) {

	log.Debugf("generate_keyset: Publishing secret IPNSKey to IPFS: %v", k.IPNSKey.PublicKey)
	err := k.IPNSKey.ExportToIPFS(*forcePublish)
	if err != nil {
		log.Warnf("create_and_print_keyset: failed to export keyset: %v", err)
	}
	log.Infof("create_and_print_keyset: exported IPNSkey to IPFS: %s", k.IPNSKey.DID)

	d, err := doc.NewFromKeyset(keyset, k.IPNSKey.DID)
	if err != nil {
		log.Fatalf("create_and_print_keyset: failed to create DOC: %v", err)
	}

	assertionMethod, err := d.GetAssertionMethod()
	if err != nil {
		log.Fatalf("create_and_print_keyset: failed to get verification method: %v", err)
	}
	d.Sign(k.SigningKey, assertionMethod)

	_, err = d.Publish()
	if err != nil {
		log.Fatalf("create_and_print_keyset: failed to publish DOC: %v", err)
	}
	log.Debugf("create_and_print_keyset: published DOC: %s", d.ID)

}

func GetKeyset() *set.Keyset {
	return keyset
}
