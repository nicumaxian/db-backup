package commands

import (
	"db-backup/configuration"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/urfave/cli/v2"
)

func ConfigurationCommands() *cli.Command {
	return &cli.Command{
		Name:  "config",
		Usage: "Manage configurations",
		Subcommands: []*cli.Command{
			configurationAddCommand(),
			configurationEditCommand(),
			configurationDeleteCommand(),
			configurationListCommand(),
		},
	}
}

func configurationAddCommand() *cli.Command {
	type SurveyAnswer struct {
		Name     string `survey:"name"`
		Host     string `survey:"host"`
		Database string `survey:"database"`
		Username string `survey:"username"`
		Password string `survey:"password"`
	}

	var qs = []*survey.Question{
		{
			Name:   "name",
			Prompt: &survey.Input{Message: "Name of database entry"},
			Validate: func(ans interface{}) error {
				str, ok := ans.(string)
				if !ok {
					return fmt.Errorf("please provide a value")
				}

				cfg, err := configuration.Read()
				if err != nil {
					return err
				}

				if _, ok := cfg.Databases[str]; ok {
					return fmt.Errorf("configuration already exists")

				}
				return nil
			},
		},
		{
			Name:   "host",
			Prompt: &survey.Input{Message: "Host", Default: "localhost"},
		},
		{
			Name:   "database",
			Prompt: &survey.Input{Message: "Database"},
		},
		{
			Name:   "username",
			Prompt: &survey.Input{Message: "Username", Default: "postgres"},
		},
		{
			Name:   "password",
			Prompt: &survey.Password{Message: "Password"},
		},
	}

	return &cli.Command{
		Name:  "add",
		Usage: "Adds a new database entry",
		Action: func(context *cli.Context) error {
			var answers SurveyAnswer
			err := survey.Ask(qs, &answers)
			if err != nil {
				return err
			}

			cfg, err := configuration.Read()
			if err != nil {
				return err
			}

			cfg.Databases[answers.Name] = configuration.DbConfiguration{
				Host:     answers.Host,
				Database: answers.Database,
				Username: answers.Username,
				Password: answers.Password,
			}

			err = configuration.Write(cfg)
			if err != nil {
				return err
			}

			fmt.Println("Data successfully saved")
			return nil
		},
	}
}

func configurationEditCommand() *cli.Command {
	return &cli.Command{
		Name:  "edit",
		Usage: "Edit an existing database configuration",
		Action: func(context *cli.Context) error {
			var name string
			err := survey.AskOne(&survey.Input{Message: "Name of database entry"}, &name, func(options *survey.AskOptions) error {
				options.Validators = append(options.Validators, validateExistingConfigEntry())
				return nil
			})
			if err != nil {
				return err
			}

			type SurveyAnswer struct {
				Host     string `survey:"host"`
				Database string `survey:"database"`
				Username string `survey:"username"`
				Password string `survey:"password"`
			}

			cfg, err := configuration.Read()
			if err != nil {
				return err
			}

			existingEntry := cfg.Databases[name]

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
				{
					Name:   "password",
					Prompt: &survey.Password{Message: "Password (skip if unchanged)"},
					Transform: func(ans interface{}) (newAns interface{}) {
						str, ok := ans.(string)
						if !ok || len(str) == 0 {
							return existingEntry.Password
						}

						return str
					},
				},
			}

			var answers SurveyAnswer
			err = survey.Ask(qs, &answers)
			if err != nil {
				return err
			}

			cfg.Databases[name] = configuration.DbConfiguration{
				Host:     answers.Host,
				Database: answers.Database,
				Username: answers.Username,
				Password: answers.Password,
			}

			err = configuration.Write(cfg)
			if err != nil {
				return err
			}

			fmt.Println("Data successfully saved")
			return nil
		},
	}
}

func validateExistingConfigEntry() func(ans interface{}) error {
	return func(ans interface{}) error {
		str, ok := ans.(string)
		if !ok {
			return fmt.Errorf("please provide a value")
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

func configurationDeleteCommand() *cli.Command {
	return &cli.Command{
		Name:  "delete",
		Usage: "Delete an existing database configuration",
		Action: func(context *cli.Context) error {
			var name string
			err := survey.AskOne(&survey.Input{Message: "Name of database entry"}, &name, func(options *survey.AskOptions) error {
				options.Validators = append(options.Validators, validateExistingConfigEntry())
				return nil
			})
			if err != nil {
				return err
			}

			var confirm bool
			prompt := &survey.Confirm{
				Message: "Are you sure you want to delete it?",
			}
			err = survey.AskOne(prompt, &confirm)
			if !confirm {
				fmt.Printf("Aborting")
				return nil
			}

			cfg, err := configuration.Read()
			if err != nil {
				return err
			}

			delete(cfg.Databases, name)

			err = configuration.Write(cfg)
			if err != nil {
				return err
			}

			fmt.Println("Data successfully saved")
			return nil
		},
	}
}

func configurationListCommand() *cli.Command {
	var name string
	return &cli.Command{
		Name:  "list",
		Usage: "Lists existing configurations",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "name",
				Destination: &name,
			},
		},
		Action: func(context *cli.Context) error {
			cfg, err := configuration.Read()
			if err != nil {
				return fmt.Errorf("failed to read configuration")
			}
			if len(name) == 0 {
				fmt.Printf("There are %v databases configured\n", len(cfg.Databases))

				for key := range cfg.Databases {
					fmt.Printf("\t%v\n", key)
				}
			} else {
				if dbConfiguration, ok := cfg.Databases[name]; ok {
					fmt.Printf("Host: %v\n",dbConfiguration.Host)
					fmt.Printf("Database: %v\n",dbConfiguration.Database)
					fmt.Printf("Username: %v\n",dbConfiguration.Username)
				}
			}

			return nil
		},
	}
}
