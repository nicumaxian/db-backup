package main

import (
	"db-backup/commands"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name:  "db-backup",
		Usage: "A tool to backup and restore database easily",
		Commands: []*cli.Command{
			commands.ConfigurationCommands(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Panic(err)
	}
}
