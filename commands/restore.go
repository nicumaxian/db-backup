package commands

import (
	"db-backup/configuration"
	"db-backup/drivers"
	"db-backup/logging"
	"db-backup/storage"
	"db-backup/utils"
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/pterm/pterm"
	"github.com/urfave/cli/v2"
	"io/fs"
	"path"
)

func RestoreCommand() *cli.Command {
	var configName string
	var latest bool
	var backupName string
	var bucket string

	latestFlag := &cli.BoolFlag{
		Name:        "latest",
		Destination: &latest,
		Usage:       "Use latest existing backup",
	}

	specificBackupFlag := &cli.StringFlag{
		Name:        "name",
		Destination: &backupName,
		Usage:       "Use specific backup",
	}

	return &cli.Command{
		Name:  "restore",
		Usage: "Restore a backup",
		Flags: []cli.Flag{
			configurationFlag(&configName),
			bucketFlag(&bucket),
			latestFlag,
			specificBackupFlag,
		},
		Action: func(context *cli.Context) error {
			err := validateFlags(latestFlag, specificBackupFlag)
			if err != nil {
				return err
			}

			err = survey.ComposeValidators(validateName(), validateExistingConfigEntry())(configName)
			if err != nil {
				return err
			}

			cfg, err := configuration.Read()
			if err != nil {
				return err
			}

			backups, directory, err := storage.GetBackups(configName, bucket)
			if err != nil {
				return err
			}

			if len(backups) == 0 {
				return errors.New("no backups found")
			}

			var usedBackup fs.FileInfo

			if latest {
				usedBackup = getLatestBackup(backups)
			}

			if len(backupName) > 0 {
				usedBackup = getBackupByName(backups, backupName)
			}

			if usedBackup == nil {
				return errors.New("please specify a backup to restore")
			}

			pterm.Printf("Restoring %s\n", pterm.Green(usedBackup.Name()))

			client, err := drivers.CreateDbClient(cfg.Databases[configName], logging.NewMockLogger())
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

func validateFlags(latestFlag *cli.BoolFlag, specificBackupFlag *cli.StringFlag) error {
	if *latestFlag.Destination && len(*specificBackupFlag.Destination) > 0 {
		return fmt.Errorf("please provide either latest or name flag")
	}

	return nil
}

func getLatestBackup(backups []fs.FileInfo) fs.FileInfo {
	usedBackup := backups[0]
	for _, el := range backups {
		if el.ModTime().After(usedBackup.ModTime()) {
			usedBackup = el
		}
	}

	return usedBackup
}

func getBackupByName(backups []fs.FileInfo, name string) fs.FileInfo {
	for _, el := range backups {
		if el.Name() == name {
			return el
		}
	}

	return nil
}
