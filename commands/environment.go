package commands

import (
	"gopkg.in/yaml.v2"
	"os"
	"path"
)

func CreateConfigurationFolderIfDoesntExist(userHomeDir string, baseDir string) error {
	configurationFolderPath := path.Join(userHomeDir, baseDir)
	if _, err := os.Stat(configurationFolderPath); os.IsNotExist(err) {
		err = os.Mkdir(configurationFolderPath, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

func CreateInitialConfigurationFileIfDoesntExist(appDir string, filename string, seed interface{}) error {
	configurationFilePath := path.Join(appDir, filename)
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
