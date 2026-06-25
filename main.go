package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yukiisen/funch/cmd"
)

var VersionCmd = &cobra.Command {
	Use: "version",
	Short: "Show Version",
	Long:  "Show App Version",
	Run: func (cmd *cobra.Command, args []string) {
		fmt.Println("funch v0.1.1")
	},
}

func main() {
	rootCmd := &cobra.Command {
		Use: "funch",
	}

	cmd.InitListCmd()
	cmd.InitRunCmd()

	rootCmd.AddCommand(VersionCmd);
	rootCmd.AddCommand(cmd.RunCmd)
	rootCmd.AddCommand(cmd.AddCmd)
	rootCmd.AddCommand(cmd.SetCmd)
	rootCmd.AddCommand(cmd.ListCmd)
	rootCmd.AddCommand(cmd.InfoCmd)
	rootCmd.AddCommand(cmd.SyncCmd)
	rootCmd.Execute()
}
