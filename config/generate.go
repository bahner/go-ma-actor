package config

import (
	"errors"
	"fmt"

	"github.com/bahner/go-ma/did/doc"
	"github.com/bahner/go-ma/key/set"
	log "github.com/sirupsen/logrus"
)

// Genreates a libp2p and actor identity and returns the keyset and the actor identity
// These are imperative, so failure to generate them is a fatal error.
func GenerateActorIdentities(name string) (string, string, error) {

	keyset_string, err := GenerateActorIdentity(name)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate actor identity: %w", err)
	}

	ni, err := GenerateNodeIdentity()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate node identity: %w", err)
	}

	return keyset_string, ni, nil
}
func GenerateActorIdentity(nick string) (string, error) {

	log.Debugf("Generating new keyset for %s", nick)
	keyset_string, err := generateKeysetString(nick)
	if err != nil {
		log.Errorf("Failed to generate new keyset: %s", err)
		return "", fmt.Errorf("failed to generate new keyset: %w", err)
	}

	// Ignore already published error. That's a good thing.
	if PublishFlag() {
		err = publishActorIdentityFromString(keyset_string)

		if err != nil {
			if errors.Is(err, doc.ErrAlreadyPublished) {
				log.Warnf("Actor document already published: %v", err)
			} else {
				return "", fmt.Errorf("failed to publish actor identity: %w", err)
			}
		}
	}

	return keyset_string, nil
}

// Generates a new keyset and returns the keyset as a string
func generateKeysetString(nick string) (string, error) {

	ks, err := set.GetOrCreate(nick)
	if err != nil {
		return "", fmt.Errorf("failed to generate new keyset: %w", err)
	}
	log.Debugf("Created new keyset: %v", ks)

	pks, err := ks.Pack()
	if err != nil {
		return "", fmt.Errorf("failed to pack keyset: %w", err)
	}
	log.Debugf("Packed keyset: %v", pks)

	return pks, nil
}

func publishActorIdentityFromString(keyset_string string) error {

	keyset, err := set.Unpack(keyset_string)
	if err != nil {
		log.Errorf("publishActorIdentityFromString: Failed to unpack keyset: %v", err)
	}

	err = PublishIdentityFromKeyset(keyset)
	if err != nil {
		return fmt.Errorf("publishActorIdentityFromString: Failed to publish keyset: %w", err)
	}

	return nil
}
