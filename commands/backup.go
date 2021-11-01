package commands

import (
	"db-backup/configuration"
	"db-backup/drivers"
	"db-backup/utils"
	"github.com/AlecAivazis/survey/v2"
	"github.com/urfave/cli/v2"
)

func BackupCommand() *cli.Command {
	var name string
	return &cli.Command{
		Name:  "backup",
		Usage: "Create a backup",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "configuration",
				Required:    true,
				Destination: &name,
			},
		},
		Action: func(context *cli.Context) error {
			err := survey.ComposeValidators(validateName(), validateExistingConfigEntry())(name)
			if err != nil {
				return err
			}

			cfg, err := configuration.Read()
			if err != nil {
				return err
			}

			path, err := utils.GetBackupPath(name)
			if err != nil {
				return err
			}

			client, err := drivers.CreateDbClient(cfg.Databases[name])
			if err != nil {
				return err
			}

			err = client.Backup(path)
			if err != nil {
				return err
			}

			return nil
		},
	}
}
