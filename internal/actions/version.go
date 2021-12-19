package actions

import (
	"errors"
	"fmt"

	cli "github.com/urfave/cli/v2"

	"github.com/lucas-inacio/avm/pkg/manager"
)

func ActionVersion(ctx *cli.Context) error {
	if ctx.NArg() != 0 {
		return errors.New("version takes no parameters")
	}

	version, err := manager.GetArduinoVersion()
	if err != nil {
		return err
	}

	release, err := manager.GetLatestRelease()
	if err != nil {
		return err
	}

	fmt.Println("Version", version, "found")
	if release.Tag != version {
		fmt.Println("Version", release.Tag, "is available. Run avm update to download it.")
	}

	return nil
}
