package actions

import (
	"context"
	"errors"
	"fmt"
	"strings"

	cli "github.com/urfave/cli/v2"

	"github.com/lucas-inacio/avm/pkg/manager"
)

func ActionInstall(ctx *cli.Context) error {
	if ctx.NArg() != 1 {
		return errors.New("provide a version to install")
	}

	alreadyInstalled := true
	dir, err := manager.GetArduinoDir()
	if err != nil {
		switch err.(type) {
		case *manager.ArduinoNotFoundError:
			alreadyInstalled = false
		default:
			return err
		}
	}

	if alreadyInstalled {
		if !ctx.Bool("force") {
			return errors.New("arduino-cli is already installed. Use avm install -f <version> to overwrite")
		}
	}

	overwrite := false
	newDir := ""
	if ctx.String("dir") != ""  && ctx.String("dir") != dir {
		newDir = ctx.String("dir")
		overwrite = true
	} else {
		newDir = dir
	}

	version := ctx.Args().First()
	fmt.Println("Downloading release " + version)
	path, err := download(newDir, version)
	if err != nil {
		return err
	}
	fmt.Println("Download completed")

	// task, err := manager.DecompressFileZip(context.Background(), path)
	// if err != nil {
	// 	return err
	// }	
	var task *manager.TaskProgress
	if strings.HasSuffix(path, ".zip") {
		task, err = manager.DecompressFileZip(context.Background(), path)
		if err != nil {
			return err
		}
	} else {
		task, err = manager.DecompressFileTargz(context.Background(), path)
		if err != nil {
			return err
		}
	}
	<- task.Done()
	fmt.Println("Installation completed")

	if overwrite {
		fmt.Println("You should update your path to reflect your changes")
		fmt.Println("New installation path:", newDir)
	}
	return nil
}