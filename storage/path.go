package storage

import (
	"os"
	"path"
)

var UserHomeDir string
var AppDir string
var BaseDir = ".db-backup"
var ConfigurationFilename = "config.yml"

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	UserHomeDir = homeDir
	AppDir = path.Join(homeDir, BaseDir)
}
