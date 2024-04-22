package config

import (
	"sync"

	"github.com/bahner/go-ma-actor/internal"
	"github.com/ipfs/boxo/keystore"
)

var (
	ks keystore.Keystore

	keystoreOnce sync.Once
)

func Keystore() keystore.Keystore {
	initKeystore()
	return ks
}

func initKeystore() {

	keystoreOnce.Do(func() {

		ks, err = keystore.NewFSKeystore(DBKeystore())
		if err != nil {
			panic(err)
		}

	})

}

func defaultKeystorePath() string {
	return internal.NormalisePath(dataHome + "/keystore")
}
