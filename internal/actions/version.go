package actions

import (
	"fmt"

	cli "github.com/urfave/cli/v2"

	"github.com/lucas-inacio/avm/pkg/manager"
)

func ActionVersion(ctx *cli.Context) error {
	version, err := manager.GetArduinoVersion()
	if err != nil {
		return err
	}

	fmt.Println(version)

	return nil
}