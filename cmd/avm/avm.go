package main

import (
	// "context"
	// "fmt"
	"log"
	"os"
	// "time"

	cli "github.com/urfave/cli/v2"

	"github.com/lucas-inacio/avm/internal/actions"
)

var CommandGet = cli.Command{
	Name: "get",
	Usage: "download arduino-cli",
	ArgsUsage: "testando argumentos",
	Action: actions.ActionGet,
}

func main() {
	cliApp := &cli.App{
		Name: "avm",
		Usage: "arduino-cli version manager",
		Authors: []*cli.Author{
			{
				Name: "Lucas Inácio Viegas",
				Email: "lucas.viegas@edu.pucrs.br",
			},
		},
		Commands: []*cli.Command{
			&CommandGet,
		},
	}

	if err := cliApp.Run(os.Args); err != nil {
		log.Fatalf("An error occurred: %v", err.Error())
	}
	// releases, err := manager.GetReleases()
	// if err != nil {
	// 	log.Fatalf("An error occurred: %v", err)
	// }

	// for _, rel := range releases {
	// 	fmt.Println(rel.Tag)
	// }
	// latest, err := manager.GetLatestRelease()
	// if err != nil {
	// 	log.Fatalf("An error occurred: %v", err.Error())
	// }
	
	// dir, err := os.Getwd()
	// if err != nil {
	// 	log.Fatalf("An error occurred: %v", err.Error())
	// }

	// ctx := context.Background()
	// fmt.Println("Downloading release " + latest.Tag)
	// task, err := manager.DownloadRelease(ctx, dir, latest.Tag)
	// if err != nil {
	// 	log.Fatalf("An error occurred: %v", err.Error())
	// }

	// ticker := time.NewTicker(time.Millisecond * 200)
	// run := true
	// for run {
	// 	select {
	// 	case <- ticker.C:
	// 		fmt.Printf("%.3v%%\n", task.GetProgress() * 100.0)
	// 	case <- task.Done():
	// 		run = false
	// 	}
	// }

	// fmt.Println("Download successful")
}