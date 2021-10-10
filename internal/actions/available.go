package actions

import (
	"errors"
	"fmt"

	cli "github.com/urfave/cli/v2"

	"github.com/lucas-inacio/avm/pkg/manager"
)

func ActionAvailable(ctx *cli.Context) error {
	if ctx.NArg() != 0 {
		return errors.New("available takes no parameters")
	}

	releases, err := manager.GetReleases()
	if err != nil {
		return err
	}

	for _, rel := range releases {
		fmt.Println(rel.Tag)
	}

	return nil
}