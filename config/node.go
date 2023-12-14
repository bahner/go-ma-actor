package config

import (
	"flag"

	"github.com/libp2p/go-libp2p/core/crypto"
	mb "github.com/multiformats/go-multibase"
	log "github.com/sirupsen/logrus"
	"go.deanishe.net/env"
)

var (
	// P2P Node identity
	nodeMultibaseIdentity = flag.String("nodeIdentity", env.Get(GO_MA_ACTOR_NODE_IDENTITY_VAR, ""),
		"Multibaseencoded libp2p privkey for the node. You can use environment variable "+GO_MA_ACTOR_NODE_IDENTITY_VAR+" to set this.")
)

func InitNodeIdentity() {

	if *nodeMultibaseIdentity == "" {
		*nodeMultibaseIdentity, _ = generateNodeIdentity()
	}

	log.Debugf("Node identity: %s", *nodeMultibaseIdentity)

}

func GetNodeMultibasePrivKey() string {

	return *nodeMultibaseIdentity
}

func GetNodeIdentity() crypto.PrivKey {

	_, privKeyBytes, err := mb.Decode(*nodeMultibaseIdentity)
	if err != nil {
		return nil
	}

	privKey, err := crypto.UnmarshalPrivateKey(privKeyBytes)
	if err != nil {
		return nil
	}

	return privKey

}

func generateNodeIdentity() (string, error) {
	pk, _, err := crypto.GenerateKeyPair(crypto.Ed25519, -1)
	if err != nil {
		log.Errorf("failed to generate node identity: %s", err)
		return "", err
	}

	pkBytes, err := pk.Raw()
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
