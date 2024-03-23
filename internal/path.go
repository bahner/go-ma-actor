package internal

import "path/filepath"

func NormalisePath(path string) string {
	return filepath.ToSlash(filepath.Clean(path))
}
