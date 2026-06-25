package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yukiisen/funch/api"
	"github.com/yukiisen/funch/config"
)

var SyncCmd = &cobra.Command {
	Use: "sync",
	Args: cobra.MaximumNArgs(0),
	RunE: sync,
}

func sync(cmd *cobra.Command, args []string) error {
	lib := config.LoadLibrary()

	sgdb := api.SGDB { Key: lib.Config.CoversKey }


	for i, g := range lib.Games {
		if len(g.Icons) == 0 {
			fmt.Printf("Syncing %s\n", g.DisplayName)

			urls, err := sgdb.GetGameIcons(g.DisplayName)
			if err != nil { return err }
			
			lib.Games[i].Icons = urls
			lib.Games[i].IconIdx = 0
		}
	}

	err := config.UpdateLibrary(lib);
	if err != nil { return err }

	return nil
}
