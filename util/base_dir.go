package util

import "os"

var (
	osWd, _ = os.Getwd()
	baseDir = os.Getenv("CWD")
)

// GetBaseDir returns the path of project directory
func GetBaseDir() string {
	if baseDir == "" {
		return osWd
	}

	return baseDir
}
