package util

import (
	"os"
	"path/filepath"
)

const uploadDir = "uploads/assets"

func GetUploadDir() string {
	return filepath.Join(GetBaseDir(), uploadDir)
}

func EnsureUploadDir() error {
	return os.MkdirAll(GetUploadDir(), os.ModePerm)
}

func GetUploadFilePath(name string) string {
	return filepath.Join(GetUploadDir(), name)
}
