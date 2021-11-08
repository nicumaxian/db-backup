package main

import (
	"db-backup/commands"
	"db-backup/configuration"
	"db-backup/storage"
	"github.com/pterm/pterm"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	app := &cli.App{
		Name:    "db-backup",
		Usage:   "A tool to backup and restore database easily",
		Version: "0.1",
		Authors: []*cli.Author{
			{
				Name:  "Nicu Maxian",
				Email: "maxiannicu@gmail.com",
			},
			{
				Name:  "Andrian Boscanean",
				Email: "boscanean.andrian@gmail.com",
			},
		},
		Commands: []*cli.Command{
			commands.BackupCommand(),
			commands.RestoreCommand(),
			commands.ConfigurationCommands(),
		},
		Before: func(ctx *cli.Context) error {
			err := storage.CreateConfigurationFolderIfDoesntExist()
			if err != nil {
				return err
			}

			err = storage.CreateInitialConfigurationFileIfDoesntExist(configuration.Default)
			if err != nil {
				return err
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		pterm.Error.Println(err)
	}
}
