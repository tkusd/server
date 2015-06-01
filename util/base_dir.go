package util

import "os"

var (
	osWd, _ = os.Getwd()
	baseDir = os.Getenv("GO_CWD")
)

// GetBaseDir returns the path of project directory
func GetBaseDir() string {
	if baseDir == "" {
		return osWd
	}

	return baseDir
}
