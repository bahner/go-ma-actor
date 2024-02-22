package config

import "os"

const defaultEditor string = "vim"

func GetEditor() string {
	editor := os.Getenv("EDITOR")

	if editor == "" {
		return defaultEditor
	}
	return editor

}
