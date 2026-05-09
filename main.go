package main

import (
	"funch/cmd"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command {
		Use: "funch",
	}

	rootCmd.AddCommand(cmd.RunCmd)
	rootCmd.AddCommand(cmd.AddCmd)
	rootCmd.Execute()
}
