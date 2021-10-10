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
		},
	}

	if err := cliApp.Run(os.Args); err != nil {
		log.Fatalf("An error occurred: %v", err.Error())
	}
}
