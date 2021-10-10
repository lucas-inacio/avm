package actions

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	cli "github.com/urfave/cli/v2"

	"github.com/lucas-inacio/avm/pkg/manager"
)

func ActionGet(cliCtx *cli.Context) error {
	latest, err := manager.GetLatestRelease()
	if err != nil {
		log.Fatalf("An error occurred: %v", err.Error())
	}
	
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("An error occurred: %v", err.Error())
	}

	ctx := context.Background()
	fmt.Println("Downloading release " + latest.Tag)
	task, err := manager.DownloadRelease(ctx, dir, latest.Tag)
	if err != nil {
		log.Fatalf("An error occurred: %v", err.Error())
	}

	ticker := time.NewTicker(time.Millisecond * 200)
	run := true
	for run {
		select {
		case <- ticker.C:
			fmt.Printf("%.3v%%\n", task.GetProgress() * 100.0)
		case <- task.Done():
			run = false
		}
	}

	fmt.Println("Download successful")

	return nil
}