package config

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

// This is a function that'll handle generateing and possible printing of the config.
// You can called his after gerenating your config to properly handle generate and show-config flags.
func Generate(c Config) {

	// Just a precaution, if the generate flag is not set, we don't want to generate the config.
	if !GenerateFlag() {
		return
	}

	if c == nil {
		log.Fatalf("No config provided.")
	}

	y, err := c.MarshalToYAML()
	if err != nil {
		log.Fatalf("Failed to marshal config to YAML: %v", err)
	}

	writeGeneratedConfigFile(y)

	os.Exit(0)
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
		panic(errMsg)
	}
	defer file.Close()

	// Write content to file.
	if _, err := file.Write(content); err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}

	log.Printf("Generated config file %s", filePath)
}
