package key

import (
	"crypto/ed25519"
	"fmt"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	mb "github.com/multiformats/go-multibase"
)

// This struct is just a simple way to pack a libp2p node key,
// in order to reuse it between runs.
// When nodes change ID they can't connect to each other easily.
type Key struct {
	ID               peer.ID
	PrivKey          crypto.PrivKey
	MultibasePrivKey string
}

func New() (*Key, error) {

	k := new(Key)
	var err error

	k.PrivKey, _, err = crypto.GenerateKeyPair(crypto.Ed25519, ed25519.PrivateKeySize)
	if err != nil {
		return nil, fmt.Errorf("key.New: failed to generate keypair: %w", err)
	}

	k.ID, err = peer.IDFromPrivateKey(k.PrivKey)
	if err != nil {
		return nil, fmt.Errorf("key.New: failed to get peer ID: %w", err)
	}

	privkeyBytes, err := k.PrivKey.Raw()
	if err != nil {
		return nil, fmt.Errorf("key.New: failed to get raw private key: %w", err)
	}

	k.MultibasePrivKey, err = mb.Encode(mb.Base58BTC, privkeyBytes)
	if err != nil {
		return nil, fmt.Errorf("key.New: failed to encode private key: %w", err)
	}

	return k, nil
}

// Extracts a libp2p node key from a multibase encoded string.
// The string is suppose to come as a command line argument.
func NewFromMultibasePrivKey(m string) (*Key, error) {

	k := new(Key)
	var err error

	k.MultibasePrivKey = m

	_, privkeyBytes, err := mb.Decode(m)
	if err != nil {
		return nil, fmt.Errorf("key.NewFromMultibasePrivKey: failed to decode private key: %w", err)
	}

	k.PrivKey, err = crypto.UnmarshalPrivateKey(privkeyBytes)
	if err != nil {
		return nil, fmt.Errorf("key.NewFromMultibasePrivKey: failed to unmarshal private key: %w", err)
	}

	k.ID, err = peer.IDFromPrivateKey(k.PrivKey)
	if err != nil {
		return nil, fmt.Errorf("key.NewFromMultibasePrivKey: failed to get peer ID: %w", err)
	}

	return k, nil
}
