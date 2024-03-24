package config

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

// This is a function that'll handle generateing and possible printing of the config.
// You can called his after gerenating your config to properly handle generate and show-config flags.
func HandleGenerate(c Config) {

	if c == nil {
		log.Fatalf("No config provided.")
	}

	y, err := c.MarshalToYAML()
	if err != nil {
		log.Fatalf("Failed to marshal config to YAML: %v", err)
	}

	writeGeneratedConfigFile(y)

	if ShowConfigFlag() {
		c.Print()
	}

	os.Exit(0)
}

// Genereates a libp2p and actor identity and returns the keyset and the actor identity
// These are imperative, so failure to generate them is a fatal error.
func generateActorIdentities(name string) (string, string, error) {

	keyset_string, err := generateActorIdentity(name)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate actor identity: %w", err)
	}

	ni, err := generateNodeIdentity()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate node identity: %w", err)
	}

	return keyset_string, ni, nil
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
