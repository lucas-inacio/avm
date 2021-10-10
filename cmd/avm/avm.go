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
	ArgsUsage: "testando argumentos",
	Action:    actions.ActionGet,
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
		},
	}

	if err := cliApp.Run(os.Args); err != nil {
		log.Fatalf("An error occurred: %v", err.Error())
	}
}
