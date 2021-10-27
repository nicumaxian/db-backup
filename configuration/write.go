package configuration

import (
	"bytes"
	"gopkg.in/yaml.v2"
)

func Write(configuration Configuration)  error {
	b, err := yaml.Marshal(&configuration)
	if err != nil {
		return err
	}

	defaultConfig := bytes.NewReader(b)

	v = initConfiguration()
	err = v.MergeConfig(defaultConfig)
	if err != nil {
		return err
	}

	err = v.WriteConfig()
	if err != nil {
		return err
	}

	return nil
}

