package storage

import (
	"gopkg.in/yaml.v2"
	"os"
	"path"
)

var userHomeDir string
var AppDir string
var baseDir = ".db-backup"
var configurationFilename = "config.yaml"

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	userHomeDir = homeDir
	AppDir = path.Join(homeDir, baseDir)
}


func CreateConfigurationFolderIfDoesntExist() error {
	configurationFolderPath := path.Join(userHomeDir, baseDir)

	_, err := os.Stat(configurationFolderPath)
	switch os.IsNotExist(err) {
	case true:
		err = os.Mkdir(configurationFolderPath, 0755)
		if err != nil {
			return err
		}
		break
	}

	return nil
}

func CreateInitialConfigurationFileIfDoesntExist(seed interface{}) error {
	configurationFilePath := path.Join(AppDir, configurationFilename)
	if _, err := os.Stat(configurationFilePath); !os.IsNotExist(err) {
		return nil
	}

	f, err := os.Create(configurationFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := yaml.Marshal(seed)
	if err != nil {
		return err
	}

	_, writeErr := f.Write(data)
	if writeErr != nil {
		return err
	}

	return nil
}
