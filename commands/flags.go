package commands

import "github.com/urfave/cli/v2"

func bucketFlag(bucket *string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "bucket",
		Destination: bucket,
		DefaultText: "generic",
		Value:       "generic",
	}
}

func configurationFlag(configuration *string) *cli.StringFlag {
	return&cli.StringFlag{
		Name:        "configuration",
		Required:    true,
		Destination: configuration,
	}
}
