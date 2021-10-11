package main

import (
	"log"
	"os"

	cli "github.com/urfave/cli/v2"

	"github.com/lucas-inacio/avm/internal/actions"
)

var CommandGet = &cli.Command{
	Name:      "get",
	Usage:     "Download arduino-cli to current directory",
	ArgsUsage: `Specify a version as returned by "avm available" to download that version`,
	Action:    actions.ActionGet,
}

var CommandVersion = &cli.Command{
	Name:      "version",
	Usage:     "Show installed arduino-cli version",
	Action:    actions.ActionVersion,
}

var CommandAvailable = &cli.Command{
	Name:      "available",
	Usage:     "Get available arduino-cli releases",
	Action:    actions.ActionAvailable,
}

var CommandUpdate = &cli.Command{
	Name:      "update",
	Usage:     "Update arduino-cli version",
	Action:    actions.ActionUpdate,
}

var CommandInstall = &cli.Command{
	Name:      "install",
	Usage:     "Install specific arduino-cli version",
	ArgsUsage: `Specify a version string as returned by "avm available"`,
	Action:    actions.ActionInstall,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name: "force",
			Aliases: []string{"f"},
			Usage: "overwrite current version",
		},
		&cli.StringFlag{
			Name: "dir",
			Aliases: []string{"d"},
			Usage: "specify full path for installation",
		},
	},
}

func main() {
	cliApp := &cli.App{
		Name:  "avm",
		Usage: "arduino-cli version manager",
		Authors: []*cli.Author{
			{
				Name:  "Lucas Inácio Viegas",
				Email: "lucas.viegas@edu.pucrs.br",
			},
		},
		Commands: []*cli.Command{
			CommandGet,
			CommandVersion,
			CommandAvailable,
			CommandUpdate,
			CommandInstall,
		},
	}

	if err := cliApp.Run(os.Args); err != nil {
		log.Fatalf("An error occurred: %v", err.Error())
	}
}
