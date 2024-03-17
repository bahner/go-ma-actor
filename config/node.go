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
