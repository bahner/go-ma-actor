package main

import (
	"flag"
	"fmt"
	"os"

	doc "github.com/bahner/go-ma/did/doc"
	keyset "github.com/bahner/go-ma/key/set"
	log "github.com/sirupsen/logrus"
)

func main() {

	fmt.Fprint(os.Stderr, "******************************************************************\n")
	fmt.Fprint(os.Stderr, "*The following strings contains secrets and should not be shared.*\n")
	fmt.Fprint(os.Stderr, "*              It is only meant for testing.                     *\n")
	fmt.Fprint(os.Stderr, "******************************************************************\n")

	log.SetLevel(log.ErrorLevel)

	name := flag.String("name", "", "Name of the entity to create")
	publish := flag.Bool("publish", false, "Publish the entity document to IPFS")
	logLevel := flag.String("loglevel", "error", "Set the log level (debug, info, warn, error, fatal, panic)")

	flag.Parse()
	_level, err := log.ParseLevel(*logLevel)
	if err != nil {
		panic(err)
	}
	log.SetLevel(_level)
	log.Debugf("main: log level set to %v", _level)

	// Create a new keyset for the entity
	keyset, err := keyset.GetOrCreate(*name)
	if err != nil {
		panic(err)
	}
	log.Debugf("main: keyset: %v", keyset)

	if *publish {
		d, err := doc.NewFromKeyset(keyset)
		if err != nil {
			panic(err)
		}

		c, err := d.Publish()
		if err != nil {
			panic(err)
		}

		log.Debugf("main: published document: %v to %v", d, c)
	}

	packedKeyset, err := keyset.Pack()
	if err != nil {
		panic(err)
	}
	fmt.Println(packedKeyset)

}
