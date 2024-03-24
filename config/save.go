package config

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

// Write the generated config to the correct file
// NB! This fails fatally in case of an error.
func Save(c Config) error {
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
		return fmt.Errorf(errMsg)
	}
	defer file.Close()

	content, err := c.MarshalToYAML()
	if err != nil {
		return fmt.Errorf("failed to marshal to YAML: %w", err)
	}

	// Write content to file.
	if _, err := file.Write(content); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	log.Printf("Generated config file %s", filePath)
	return nil
}
