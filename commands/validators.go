package commands

import (
	"db-backup/configuration"
	"fmt"
)

func validateExistingConfigEntry() func(ans interface{}) error {
	return func(ans interface{}) error {
		str, ok := ans.(string)
		if !ok {
			return fmt.Errorf("please provide a configuration name")
		}

		cfg, err := configuration.Read()
		if err != nil {
			return err
		}

		if _, ok := cfg.Databases[str]; !ok {
			return fmt.Errorf("configuration does not exist")

		}
		return nil
	}
}


func validateNotExistingConfigEntry() func(ans interface{}) error {
	return func(ans interface{}) error {
		str, ok := ans.(string)
		if !ok {
			return fmt.Errorf("please provide a configuration name")
		}

		cfg, err := configuration.Read()
		if err != nil {
			return err
		}

		if _, ok := cfg.Databases[str]; ok {
			return fmt.Errorf("configuration already exist")

		}
		return nil
	}
}