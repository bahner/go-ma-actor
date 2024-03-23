package config

import (
	"github.com/libp2p/go-libp2p/core/crypto"
	mb "github.com/multiformats/go-multibase"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NodeIdentity() crypto.PrivKey {

	log.Debugf("config.GetNodeIdentity: %s", viper.GetString("p2p.identity"))
	_, privKeyBytes, err := mb.Decode(viper.GetString("p2p.identity"))
	if err != nil {
		log.Debugf("config.GetNodeIdentity: Failed to decode node identity: %v", err)
		return nil
	}

	privKey, err := crypto.UnmarshalPrivateKey(privKeyBytes)
	if err != nil {
		log.Debugf("config.GetNodeIdentity: Failed to unmarshal node identity: %v", err)
		return nil
	}

	log.Debug("Config.GetNodeIdentity: ", privKey.GetPublic())
	return privKey

}

func GenerateNodeIdentity() (string, error) {
	pk, _, err := crypto.GenerateKeyPair(crypto.Ed25519, -1)
	if err != nil {
		log.Errorf("failed to generate node identity: %s", err)
		return "", err
	}

	pkBytes, err := crypto.MarshalPrivateKey(pk)
	if err != nil {
		log.Errorf("failed to generate node identity: %s", err)
		return "", err
	}

	ni, err := mb.Encode(mb.Base58BTC, pkBytes)
	if err != nil {
		log.Errorf("failed to encode node identity: %s", err)
		return "", err
	}

	return ni, nil

}
