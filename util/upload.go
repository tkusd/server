package util

import (
	"os"
	"path/filepath"

	"github.com/tkusd/server/config"
)

func GetAssetFilePath(name string) string {
	return filepath.Join(config.BaseDir, config.Config.AssetDir, name)
}

func IsAssetExist(name string) bool {
	path := GetAssetFilePath(name)

	if _, err := os.Stat(path); err == nil {
		return true
	}

	return false
}
