package commands

import (
	"db-backup/configuration"
	"github.com/AlecAivazis/survey/v2"
)

func promptValidConfigurationName() (string, error) {
	var name string
	err := survey.AskOne(&survey.Input{Message: "Name of database entry"}, &name, survey.WithValidator(validateExistingConfigEntry()))
	if err != nil {
		return "", err
	}

	return name, nil
}

func promptFreeConfigurationName() (string, error) {
	var name string
	err := survey.AskOne(&survey.Input{Message: "Name of database entry"}, &name, survey.WithValidator(validateNotExistingConfigEntry()))
	if err != nil {
		return "", err
	}

	return name, nil
}

func promptConfigurationEntry(existingEntry configuration.DbConfiguration) (configuration.DbConfiguration, error) {
	password := &survey.Question{
		Name:   "password",
		Prompt: &survey.Password{Message: "Password (skip if unchanged)"},
		Transform: func(ans interface{}) (newAns interface{}) {
			str, ok := ans.(string)
			if !ok || len(str) == 0 {
				return existingEntry.Password
			}

			return str
		},
	}
	if existingEntry.Password == "" {
		password = &survey.Question{
			Name:   "password",
			Prompt: &survey.Password{Message: "Password"},
		}
	}
	var qs = []*survey.Question{
		{
			Name:   "host",
			Prompt: &survey.Input{Message: "Host", Default: existingEntry.Host},
		},
		{
			Name:   "database",
			Prompt: &survey.Input{Message: "Database", Default: existingEntry.Database},
		},
		{
			Name:   "username",
			Prompt: &survey.Input{Message: "Username", Default: existingEntry.Username},
		},
		password,
	}

	type SurveyAnswer struct {
		Host     string `survey:"host"`
		Database string `survey:"database"`
		Username string `survey:"username"`
		Password string `survey:"password"`
	}

	var answers SurveyAnswer
	err := survey.Ask(qs, &answers)
	if err != nil {
		return configuration.DbConfiguration{}, err
	}

	return configuration.DbConfiguration{
		Host:     answers.Host,
		Database: answers.Database,
		Username: answers.Username,
		Password: answers.Password,
	}, nil
}
