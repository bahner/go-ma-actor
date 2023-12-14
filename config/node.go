package config

import (
	"flag"

	"github.com/libp2p/go-libp2p/core/crypto"
	mb "github.com/multiformats/go-multibase"
	"go.deanishe.net/env"
)

var (
	// P2P Node identity
	nodeMultibasePrivKey = flag.String("nodeKeyset", env.Get(GO_MA_ACTOR_NODE_IDENTITY_VAR, ""),
		"Multibaseencoded libp2p privkey for the node. You can use environment variable "+GO_MA_ACTOR_NODE_IDENTITY_VAR+" to set this.")
)

func GetNodeMultibasePrivKey() string {

	return *nodeMultibasePrivKey
}

func GetNodeIdentity() crypto.PrivKey {

	_, privKeyBytes, err := mb.Decode(*nodeMultibasePrivKey)
	if err != nil {
		return nil
	}

	privKey, err := crypto.UnmarshalPrivateKey(privKeyBytes)
	if err != nil {
		return nil
	}

	return privKey

}
