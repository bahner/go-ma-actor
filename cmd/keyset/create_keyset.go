package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"os"

	doc "github.com/bahner/go-ma/did/doc"
	"github.com/bahner/go-ma/key/set"
	"github.com/libp2p/go-libp2p/core/crypto"
	log "github.com/sirupsen/logrus"
)

func main() {

	fmt.Fprint(os.Stderr, "******************************************************************\n")
	fmt.Fprint(os.Stderr, "*The following strings contains secrets and should not be shared.*\n")
	fmt.Fprint(os.Stderr, "*              It is only meant for testing.                     *\n")
	fmt.Fprint(os.Stderr, "******************************************************************\n")

	log.SetLevel(log.ErrorLevel)

	name := flag.String("name", "", "(Nick)name of the entity to create")
	publish := flag.Bool("publish", false, "Publish the entity document to IPFS")
	logLevel := flag.String("loglevel", "error", "Set the log level (debug, info, warn, error, fatal, panic)")

	flag.Parse()
	_level, err := log.ParseLevel(*logLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(_level)
	log.Debugf("main: log level set to %v", _level)

	// Create a new keyset for the entity
	privKey, _, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	keyset, err := set.New(privKey, *name)
	if err != nil {
		log.Fatal(err)
	}
	log.Debugf("main: keyset: %v", keyset)

	if *publish {
		d, err := doc.NewFromKeyset(keyset)
		if err != nil {
			log.Fatal(err)
		}

		c, err := d.Publish()
		if err != nil {
			log.Fatal(err)
		}

		log.Debugf("main: published document: %v to %v", d, c)
	}

	packedKeyset, err := keyset.Pack()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(packedKeyset)

}
