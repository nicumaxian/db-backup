package commands

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

var homedir string

func init() {
	dirname, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	homedir = dirname
}

func CreateConfigurationFolderIfDoesntExist(relativePath string) error {
	configurationFolderPath := fmt.Sprintf("%s/%s", homedir, relativePath)
	if _, err := os.Stat(configurationFolderPath); os.IsNotExist(err) {
		err = os.Mkdir(configurationFolderPath, 0755)
		if err != nil {
			log.Fatal(err, "Failed to create configuration dir.")
			return err
		}
	}

	return nil
}

func CreateInitialConfigurationFileIfDoesntExist(dir string, filename string, seed interface{}) error {
	configurationFile := fmt.Sprintf("%s/%s/%s", homedir, dir, filename)
	if _, err := os.Stat(configurationFile); !os.IsNotExist(err) {
		return nil
	}

	f, err := os.Create(configurationFile)
	if err != nil {
		log.Fatal(err, "Failed to create initial configuration file.")
		return err
	}
	defer f.Close()

	data, err := yaml.Marshal(seed)
	if err != nil {
		log.Fatal(err, "Failed to serialize seed data.")
		return err
	}

	_, writeErr := f.Write(data)
	if writeErr != nil {
		log.Fatal(err, "Failed to write seed data.")
		return err
	}

	return nil
}
