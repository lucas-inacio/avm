package main

import (
	"log"
	"os"

	cli "github.com/urfave/cli/v2"

	"github.com/lucas-inacio/avm/internal/actions"
)

var CommandGet = &cli.Command{
	Name:      "get",
	Usage:     "download arduino-cli",
	ArgsUsage: `specify a version as returned by "avm version" to download that version`,
	Action:    actions.ActionGet,
}

var CommandVersion = &cli.Command{
	Name:      "version",
	Usage:     "shows installed arduino-cli version",
	Action:    actions.ActionVersion,
}

var CommandAvailable = &cli.Command{
	Name:      "available",
	Usage:     "lists available arduino-cli releases",
	Action:    actions.ActionAvailable,
}

var CommandUpdate = &cli.Command{
	Name:      "update",
	Usage:     "update arduino-cli version",
	Action:    actions.ActionUpdate,
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
		},
	}

	if err := cliApp.Run(os.Args); err != nil {
		log.Fatalf("An error occurred: %v", err.Error())
	}
}
