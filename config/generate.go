package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/bahner/go-ma/did/doc"
	"github.com/bahner/go-ma/key/set"
	"github.com/libp2p/go-libp2p/core/crypto"
	mb "github.com/multiformats/go-multibase"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func Generate(configMap map[string]interface{}) {

	if configMap == nil {
		log.Fatalf("No template set.")
	}

	// Convert the config map to YAML
	configYAML, err := yaml.Marshal(configMap)
	if err != nil {
		log.Fatalf("Failed to marshal config to YAML: %v", err)
	}

	if GenerateFlag() {
		writeGeneratedConfigFile(configYAML)
	} else {
		fmt.Println(string(configYAML))
	}
}

// Write the generated config to the correct file
// NB! This fails fatally in case of an error.
func writeGeneratedConfigFile(content []byte) {
	filePath := File()
	var errMsg string

	// Determine the file open flags based on the forceFlag
	var flags int
	if ForceFlag() {
		// Allow overwrite
		log.Warnf("Force flag set, overwriting existing config file %s", filePath)
		flags = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	} else {
		// Prevent overwrite
		flags = os.O_WRONLY | os.O_CREATE | os.O_EXCL
	}

	file, err := os.OpenFile(filePath, flags, configFileMode)
	if err != nil {
		if os.IsExist(err) {
			errMsg = fmt.Sprintf("File %s already exists.", filePath)
		} else {
			errMsg = fmt.Sprintf("Failed to open file: %v", err)
		}
		log.Fatalf(errMsg)
	}
	defer file.Close()

	// Write content to file.
	if _, err := file.Write(content); err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}

	log.Printf("Generated config file %s", filePath)
}

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
