package config

import "os"

func createXDGDirectories() error {

	err := os.MkdirAll(configHome, configDirMode)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dataHome, dataHomeMode)
	if err != nil {
		return err
	}

	return nil

}
