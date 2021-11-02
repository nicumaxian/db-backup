package commands

import (
	"db-backup/configuration"
	"fmt"
	"regexp"
)

var nameRegex = regexp.MustCompile(`^[A-z\d\-\_]+$`)
var portRegex = regexp.MustCompile("^[0-9]+$")

func validateExistingConfigEntry() func(ans interface{}) error {
	return func(ans interface{}) error {
		str, _ := ans.(string)

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
		str, _ := ans.(string)

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

func validateName() func(ans interface{}) error {
	return func(ans interface{}) error {
		str, ok := ans.(string)
		if !ok {
			return fmt.Errorf("please provide a configuration name")
		}

		if !nameRegex.MatchString(str) {
			return fmt.Errorf("configuration name can contain only characters, digits, hyphens or underscore")
		}

		return nil
	}
}

func validatePort() func(ans interface{}) error {
	return func(ans interface{}) error {
		str, ok := ans.(string)
		if !ok {
			return fmt.Errorf("please provide port number")
		}

		if !portRegex.MatchString(str) {
			return fmt.Errorf("port number should contain digits only")
		}

		port, _ := ans.(int32)
		if port > 65536 {
			return fmt.Errorf("port number cannot be greater than 65536")
		}

		return nil
	}
}
