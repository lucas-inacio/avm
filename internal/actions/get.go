package actions

import (
	"context"
	"fmt"
	"os"
	"time"

	cli "github.com/urfave/cli/v2"

	"github.com/gosuri/uiprogress"

	"github.com/lucas-inacio/avm/pkg/manager"
)

func ActionGet(cliCtx *cli.Context) error {
	latest, err := manager.GetLatestRelease()
	if err != nil {
		return err
	}

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	ctx := context.Background()
	fmt.Println("Downloading release " + latest.Tag)
	task, err := manager.DownloadRelease(ctx, dir, latest.Tag)
	if err != nil {
		return err
	}

	uiprogress.Start()
	bar := uiprogress.AddBar(int(task.GetTotal()))
	bar.AppendCompleted()

	ticker := time.NewTicker(time.Millisecond * 200)
	run := true
	for run {
		select {
		case <-ticker.C:
			bar.Set(task.GetProgress())
		case <-task.Done():
			run = false
		}
	}

	// Allow the progress bar to reach 100%
	bar.Set(task.GetProgress())
	time.Sleep(time.Second)
	fmt.Println("Download successful")

	return nil
}
