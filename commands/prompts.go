package commands

import (
	"db-backup/configuration"
	"db-backup/utils"
	"github.com/AlecAivazis/survey/v2"
)

func promptValidConfigurationName() (string, error) {
	var name string
	err := survey.AskOne(&survey.Input{Message: "Configuration name"}, &name, survey.WithValidator(survey.ComposeValidators(validateName(), validateExistingConfigEntry())))
	if err != nil {
		return "", err
	}

	return name, nil
}

func promptFreeConfigurationName() (string, error) {
	var name string
	err := survey.AskOne(&survey.Input{Message: "Configuration name"}, &name, survey.WithValidator(survey.ComposeValidators(validateName(), validateNotExistingConfigEntry())))
	if err != nil {
		return "", err
	}

	return name, nil
}

func getPortSuggestionByDriver(toComplete string, driver string) []string {
	switch driver {
	case "postgres":
		return []string{"5432", "5433"}
	case "mysql":
		return []string{"3306"}
	}

	return []string{}
}

func promptConfigurationEntry(existingEntry configuration.DbConfiguration) (configuration.DbConfiguration, error) {
	type SurveyAnswer struct {
		Driver   string `survey:"driver"`
		Host     string `survey:"host"`
		Port     string `survey:"port"`
		Database string `survey:"database"`
		Username string `survey:"username"`
		Password string `survey:"password"`
	}
	var answers SurveyAnswer

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
			Name: "driver",
			Prompt: &survey.Select{
				Message: "Driver",
				Default: utils.StrCoalesce(existingEntry.Driver, "postgres"),
				Options: []string{"postgres", "mysql"},
			},
		},
		{
			Name: "host",
			Prompt: &survey.Input{
				Message: "Host",
				Default: utils.StrCoalesce(existingEntry.Host, "localhost"),
			},
		},
		{
			Name: "port",
			Prompt: &survey.Input{
				Message: "Port",
				Default: existingEntry.Port,
				Suggest: func(toComplete string) []string {
					return getPortSuggestionByDriver(toComplete, answers.Driver)
				},
			},
			Validate: validatePort(),
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

	err := survey.Ask(qs, &answers)
	if err != nil {
		return configuration.DbConfiguration{}, err
	}

	return configuration.DbConfiguration{
		Driver:   answers.Driver,
		Host:     answers.Host,
		Port:     answers.Port,
		Database: answers.Database,
		Username: answers.Username,
		Password: answers.Password,
	}, nil
}
