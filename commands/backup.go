package commands

import (
	"db-backup/configuration"
	"db-backup/drivers"
	"db-backup/storage"
	"db-backup/utils"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/briandowns/spinner"
	"github.com/pterm/pterm"
	"github.com/urfave/cli/v2"
	"os"
	"path"
	"time"
)

func BackupCommand() *cli.Command {
	return &cli.Command{
		Name:  "backup",
		Usage: "Create a backup",
		Subcommands: []*cli.Command{
			backupCreateCommand(),
			backupListCommand(),
			backupDeleteCommand(),
			backupCleanCommand(),
		},
	}
}

func backupCreateCommand() *cli.Command {
	var name string
	var bucket string
	return &cli.Command{
		Name: "create",
		Flags: []cli.Flag{
			configurationFlag(&name),
			bucketFlag(&bucket),
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

			path, err := storage.GetNewBackupPath(name, bucket)
			if err != nil {
				return err
			}

			client, err := drivers.CreateDbClient(cfg.Databases[name])
			if err != nil {
				return err
			}

			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
			s.Start()
			err = client.Backup(path)
			s.Stop()
			if err != nil {
				return err
			}

			pterm.Println(path)
			return nil
		},
	}
}

func backupListCommand() *cli.Command {
	var name string
	var bucket string
	return &cli.Command{
		Name: "list",
		Flags: []cli.Flag{
			configurationFlag(&name),
			bucketFlag(&bucket),
		},
		Action: func(context *cli.Context) error {
			err := survey.ComposeValidators(validateName(), validateExistingConfigEntry())(name)
			if err != nil {
				return err
			}

			result, _, err := storage.GetBackups(name, bucket)
			if err != nil {
				return err
			}

			if len(result) == 0 {
				pterm.Println("There are no backups")
				return nil
			}

			data := pterm.TableData{
				{"Name", "Creation Date", "Size"},
			}

			for _, el := range result {
				data = append(data, []string{
					el.Name(),
					el.ModTime().Format(time.RFC822),
					utils.ByteCountSI(el.Size()),
				})
			}
			pterm.DefaultTable.WithHasHeader().WithData(data).Render()

			return nil
		},
	}
}

func backupDeleteCommand() *cli.Command {
	var names = &cli.StringSlice{}
	var config string
	var bucket string
	return &cli.Command{
		Name:        "delete",
		Description: "Delete existing backup file(s)",
		Flags: []cli.Flag{
			configurationFlag(&config),
			bucketFlag(&bucket),
			&cli.StringSliceFlag{
				Name:        "name",
				Required:    true,
				Destination: names,
			},
		},
		Action: func(context *cli.Context) error {
			err := survey.ComposeValidators(validateExistingConfigEntry())(config)
			if err != nil {
				return err
			}

			backups, location, err := storage.GetBackups(config, bucket)
			if err != nil {
				return err
			}

			var backupsToDelete []os.FileInfo

			for _, name := range names.Value() {
				specificBackup := utils.GetFileByName(backups, name)
				if specificBackup == nil {
					return fmt.Errorf("no such file - %s", name)
				}
				backupsToDelete = append(backupsToDelete, specificBackup)
			}

			for _, file := range backupsToDelete {
				s, _ := pterm.DefaultSpinner.Start("Deleting..")
				fullPath := path.Join(location, file.Name())
				err = os.Remove(fullPath)
				if err != nil {
					s.Fail()
					return err
				}
				s.Success("Deleted ", fullPath)
			}

			return nil
		},
	}
}

func backupCleanCommand() *cli.Command {
	var config string
	var bucket string
	var duration time.Duration
	return &cli.Command{
		Name:        "clean",
		Description: "Delete old backups",
		Flags: []cli.Flag{
			configurationFlag(&config),
			bucketFlag(&bucket),
			&cli.DurationFlag{
				Name: "duration",
				Required: true,
				Destination: &duration,
				Usage: "clean backups older than",
			},
		},
		Action: func(context *cli.Context) error {
			err := survey.ComposeValidators(validateExistingConfigEntry())(config)
			if err != nil {
				return err
			}

			backups, location, err := storage.GetBackups(config, bucket)
			if err != nil {
				return err
			}

			var backupsToDelete []os.FileInfo

			var now = time.Now()
			for _, el := range backups {
				expirationTIme := el.ModTime().Add(duration)

				if expirationTIme.Before(now) {
					backupsToDelete = append(backupsToDelete, el)
				}
			}


			for _, file := range backupsToDelete {
				s, _ := pterm.DefaultSpinner.Start("Deleting..")
				fullPath := path.Join(location, file.Name())
				err = os.Remove(fullPath)
				if err != nil {
					s.Fail()
					return err
				}
				s.Success("Deleted ", fullPath)
			}

			return nil
		},
	}
}
