package main

import (
	"fmt"
	"os"

	"github.com/bahner/go-ma/key/set"
	log "github.com/sirupsen/logrus"
)

func generateKeyset(name string, publish bool) {

	if keyset != "" {
		log.Fatal("You can't set a secret key and generate a new one at the same time.")
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
		fmt.Println("export GO_MA_ACTOR_IDENTITY=" + pks)
	} else {
		fmt.Println(pks)
	}

	log.Debugf("generate_keyset: Generated new keyset: %v", ks)

	if publish {
		log.Debugf("generate_keyset: Publishing secret IPNSKey to IPFS: %v", ks.IPNSKey.PublicKey)
		err = ks.IPNSKey.ExportToIPFS(name, *forcePublish)
		if err != nil {
			log.Fatalf("create_and_print_keyset: failed to export keyset: %v", err)
		}
		log.Infof("create_and_print_keyset: exported IPNSkey to IPFS: %s", ks.IPNSKey.DID)
	}

	os.Exit(0)
}
