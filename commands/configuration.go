package commands

import (
	"db-backup/configuration"
	"db-backup/drivers"
	"db-backup/logging"
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
			configurationTestCommand(),
		},
	}
}

func configurationAddCommand() *cli.Command {
	return &cli.Command{
		Name:  "add",
		Usage: "Adds a new database entry",
		Action: func(context *cli.Context) error {
			name, err := promptFreeConfigurationName()
			if err != nil {
				return err
			}

			entry, err := promptConfigurationEntry(configuration.DbConfiguration{})
			if err != nil {
				return err
			}

			cfg, err := configuration.Read()
			if err != nil {
				return err
			}

			cfg.Databases[name] = entry

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
			name, err := promptValidConfigurationName()
			if err != nil {
				return err
			}

			cfg, err := configuration.Read()
			if err != nil {
				return err
			}

			existingEntry := cfg.Databases[name]

			existingEntry, err = promptConfigurationEntry(existingEntry)
			if err != nil {
				return err
			}

			cfg.Databases[name] = existingEntry

			err = configuration.Write(cfg)
			if err != nil {
				return err
			}

			fmt.Println("Data successfully saved")
			return nil
		},
	}
}

func configurationTestCommand() *cli.Command {
	var verbose bool
	return &cli.Command{
		Name:  "test",
		Usage: "Test an existing database configuration",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "verbose",
				Destination: &verbose,
			},
		},
		Action: func(context *cli.Context) error {
			name, err := promptValidConfigurationName()
			if err != nil {
				return err
			}

			cfg, err := configuration.Read()
			if err != nil {
				return err
			}

			logger := logging.NewMockLogger()
			if verbose {
				logger = logging.NewBuiltinLogger()
			}

			testingConfig := cfg.Databases[name]
			client, err := drivers.CreateDbClient(testingConfig, logger)
			if err != nil {
				return err
			}

			err = client.TestConnection()
			if err != nil {
				return err
			}

			return nil
		},
	}
}

func configurationDeleteCommand() *cli.Command {
	return &cli.Command{
		Name:  "delete",
		Usage: "Delete an existing database configuration",
		Action: func(context *cli.Context) error {
			name, err := promptValidConfigurationName()
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
	return &cli.Command{
		Name:  "list",
		Usage: "Lists existing configurations",
		Action: func(context *cli.Context) error {
			cfg, err := configuration.Read()

			if err != nil {
				return fmt.Errorf("failed to read configuration")
			}
			fmt.Printf("There are %v configurations\n", len(cfg.Databases))

			for key := range cfg.Databases {
				fmt.Printf("\t%v\n", key)
			}

			return nil
		},
	}
}
