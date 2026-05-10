package main

import (
	"funch/cmd"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command {
		Use: "funch",
	}

	cmd.InitListCmd()
	cmd.InitRunCmd()

	rootCmd.AddCommand(cmd.RunCmd)
	rootCmd.AddCommand(cmd.AddCmd)
	rootCmd.AddCommand(cmd.SetCmd)
	rootCmd.AddCommand(cmd.ListCmd)
	rootCmd.AddCommand(cmd.InfoCmd)
	rootCmd.Execute()
}
