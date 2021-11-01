package configuration

import "github.com/spf13/viper"

type NotFoundErr struct {}

func (n NotFoundErr) Error() string {
	return "config file not found"
}

func Read() (Configuration, error) {
	configuration := Configuration{}

	err := v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			//create empty config file

			return Configuration{
				Databases: map[string]DbConfiguration{},
			}, nil
		}

		return configuration, err
	}

	err = v.Unmarshal(&configuration)
	if err != nil {
		return configuration, err
	}

	return configuration, nil
}

