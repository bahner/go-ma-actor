package db

import (
	"github.com/bahner/go-ma-actor/config"
	"github.com/libp2p/go-libp2p/core/crypto"
	log "github.com/sirupsen/logrus"
)

func GetOrCreateIdentity(name string) (crypto.PrivKey, error) {

	keystore := config.Keystore()

	hasKey, err := keystore.Has(name)
	if hasKey {
		return keystore.Get(name)
	}
	if err != nil {
		log.Errorf("failed to get private key from keystore: %s", err)
		return nil, err
	}

	privKey, _, err := crypto.GenerateKeyPair(crypto.Ed25519, -1)
	if err != nil {
		log.Errorf("failed to generate node identity: %s", err)
		return nil, err
	}

	err = keystore.Put(name, privKey)
	if err != nil {
		log.Errorf("failed to store private key to keystore: %s", err)
		return nil, err
	}

	return privKey, nil

}
