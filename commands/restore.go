package commands

import (
	"db-backup/configuration"
	"db-backup/drivers"
	"db-backup/storage"
	"db-backup/utils"
	"errors"
	"github.com/AlecAivazis/survey/v2"
	"github.com/pterm/pterm"
	"github.com/urfave/cli/v2"
	"io/fs"
	"path"
)

func RestoreCommand() *cli.Command {
	var name string
	var latest bool
	return &cli.Command{
		Name:  "restore",
		Usage: "Restore a backup",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "configuration",
				Required:    true,
				Destination: &name,
			},
			&cli.BoolFlag{
				Name:        "latest",
				Destination: &latest,
				Usage:       "Use latest existing backup",
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

			backups, directory, err := storage.GetBackups(name)
			if err != nil {
				return err
			}

			if len(backups) == 0 {
				return errors.New("no backups found")
			}

			var usedBackup fs.FileInfo

			if latest {
				usedBackup = backups[0]

				for _, el := range backups {
					if el.ModTime().After(usedBackup.ModTime()) {
						usedBackup = el
					}
				}
			}

			if usedBackup == nil {
				return errors.New("please specify a backup to restore")
			}

			pterm.Printf("Restoring %s\n", pterm.Green(usedBackup.Name()))

			client, err := drivers.CreateDbClient(cfg.Databases[name])
			if err != nil {
				return err
			}

			pterm.Printf("Restoring %v - %v", usedBackup.Name(), utils.ByteCountSI(usedBackup.Size()))

			err = client.Restore(path.Join(directory, usedBackup.Name()))
			if err != nil {
				return err
			}

			return nil
		},
	}
}
