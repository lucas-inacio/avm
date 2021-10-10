package actions

import (
	"context"
	"errors"
	"fmt"
	"time"

	cli "github.com/urfave/cli/v2"

	"github.com/gosuri/uiprogress"

	"github.com/lucas-inacio/avm/pkg/manager"
)

func shouldUpdate() (string, error) {
	installedVersion, err := manager.GetArduinoVersion()
	if err != nil {
		switch err.(type) {
		case *manager.ArduinoNotFoundError:
			fmt.Println("You currently don't have arduino-cli installed. Use avm install.")
		}
		return "", err
	}

	release, err := manager.GetLatestRelease()
	if err != nil {
		return "", err
	}

	if installedVersion == release.Tag {
		return "", nil
	} else {
		return release.Tag, nil
	}
}

func download(dir, version string) (string, error) {
	ctx := context.Background()
	task, err := manager.DownloadRelease(ctx, dir, version)
	if err != nil {
		return "", err
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

	bar.Set(task.GetProgress())
	time.Sleep(time.Second)

	return task.GetData().(string), nil
}

func ActionUpdate(cliCtx *cli.Context) error {
	if cliCtx.NArg() != 0 {
		return errors.New("update takes no parameters")
	}

	version, err := shouldUpdate()
	if err != nil {
		return err
	}

	if version == "" {
		fmt.Println("You already have the latest version.")
		return nil
	}

	dir, err := manager.GetArduinoDir()
	if err != nil {
		return err
	}

	fmt.Println("Downloading release " + version)
	path, err := download(dir, version)
	if err != nil {
		return err
	}
	fmt.Println("Download completed")
	
	task, err := manager.DecompressFileZip(context.Background(), path)
	if err != nil {
		return err
	}

	<- task.Done()
	fmt.Println("Installation completed")
	return nil
}