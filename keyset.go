package main

import (
	"fmt"
	"os"

	"github.com/bahner/go-ma/did/doc"
	"github.com/bahner/go-ma/key/set"
	log "github.com/sirupsen/logrus"
)

func generateKeyset(name string) {

	if keyset != "" {
		log.Info("Ignoring keyset flag as generate is set.")
	}

	if nick == "ghost" {
		log.Fatal("You need to set a nick when generating an identity.")
	}

	ks, err := set.New(name)
	if err != nil {
		log.Fatalf("Failed to generate new keyset: %v", err)
	}

	pks, err := ks.Pack()
	if err != nil {
		log.Fatalf("Failed to pack keyset: %v", err)
	}

	if *genenv {
		fmt.Println("export GO_MA_ACTOR_KEYSET=" + pks)
	} else {
		fmt.Println(pks)
	}

	log.Debugf("generate_keyset: Generated new keyset: %v", ks)

	if *publish || *forcePublish {
		log.Debugf("generate_keyset: Publishing secret IPNSKey to IPFS: %v", ks.IPNSKey.PublicKey)
		err = ks.IPNSKey.ExportToIPFS(name, *forcePublish)
		if err != nil {
			log.Fatalf("create_and_print_keyset: failed to export keyset: %v", err)
		}
		log.Infof("create_and_print_keyset: exported IPNSkey to IPFS: %s", ks.IPNSKey.DID)

		d, err := doc.NewFromKeyset(&ks, ks.IPNSKey.DID)
		if err != nil {
			log.Fatalf("create_and_print_keyset: failed to create DOC: %v", err)
		}

		assertionMethod, err := d.GetAssertionMethod()
		if err != nil {
			log.Fatalf("create_and_print_keyset: failed to get verification method: %v", err)
		}
		d.Sign(ks.SigningKey, assertionMethod)

		_, err = d.Publish()
		if err != nil {
			log.Fatalf("create_and_print_keyset: failed to publish DOC: %v", err)
		}
		log.Debugf("create_and_print_keyset: published DOC: %s", d.ID)
	}

	os.Exit(0)
}
