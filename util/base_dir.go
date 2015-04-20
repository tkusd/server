package util

import "os"

var (
	osWd, _ = os.Getwd()
	baseDir = os.Getenv("CWD")
)

func GetBaseDir() string {
	if baseDir == "" {
		return osWd
	} else {
		return baseDir
	}
}
