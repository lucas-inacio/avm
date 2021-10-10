package actions

import (
	"fmt"

	cli "github.com/urfave/cli/v2"

	"github.com/lucas-inacio/avm/pkg/manager"
)

func ActionAvailable(ctx *cli.Context) error {
	releases, err := manager.GetReleases()
	if err != nil {
		return err
	}

	for _, rel := range releases {
		fmt.Println(rel.Tag)
	}

	return nil
}