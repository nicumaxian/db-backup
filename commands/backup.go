package commands

import (
	"db-backup/configuration"
	"db-backup/drivers"
	"db-backup/storage"
	"db-backup/utils"
	"github.com/AlecAivazis/survey/v2"
	"github.com/briandowns/spinner"
	"github.com/pterm/pterm"
	"github.com/urfave/cli/v2"
	"time"
)

func BackupCommand() *cli.Command {
	return &cli.Command{
		Name:  "backup",
		Usage: "Create a backup",
		Subcommands: []*cli.Command{
			backupCreateCommand(),
			backupListCommand(),
		},
	}
}

func backupCreateCommand() *cli.Command {
	var name string
	return &cli.Command{
		Name: "create",
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

			path, err := storage.GetNewBackupPath(name)
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
	return &cli.Command{
		Name: "list",
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

			result, err := storage.GetBackups(name)
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
