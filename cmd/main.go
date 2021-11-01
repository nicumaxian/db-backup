package main

import (
	"db-backup/commands"
	"db-backup/configuration"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name:    "db-backup",
		Usage:   "A tool to backup and restore database easily",
		Version: "0.1",
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "Nicu Maxian",
				Email: "maxiannicu@gmail.com",
			},
		},
		Commands: []*cli.Command{
			commands.ConfigurationCommands(),
		},
		Before: func(ctx *cli.Context) error {
			commands.CreateConfigurationFolderIfDoesntExist(".db-backup")
			commands.CreateInitialConfigurationFileIfDoesntExist(".db-backup", "config.yml",
				configuration.Configuration{Databases: map[string]configuration.DbConfiguration{}})

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Panic(err)
	}
}
