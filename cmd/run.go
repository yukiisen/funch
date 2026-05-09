package cmd

import (
	// "funch/config"
	"github.com/spf13/cobra"
)

var RunCmd = &cobra.Command {
	Use: "run <game>",
	Run: run,
}

func run(cmd *cobra.Command, args []string) {

}
