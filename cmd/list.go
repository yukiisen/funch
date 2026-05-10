package cmd

import (
	"os"
	"fmt"
	"funch/config"
	"text/tabwriter"
	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command {
	Use: "list",
	Short: "List all games",
	Long:  "Shows all games currently stored in the library, including those imported from Steam and Lutris.",
	RunE: list,
}

func InitListCmd() {
	ListCmd.Flags().Bool("flat", false, "Display unformatted names only")
}

func list(cmd *cobra.Command, args []string) error {
	lib := config.LoadLibrary()

	flat, err := cmd.Flags().GetBool("flat")
	if err != nil { return err }

	// flat output
	if flat {
		names, _ := flatGames(cmd, args, "")

		for _, name := range names {
			fmt.Println(name)
		}

		return nil
	}

	// do formatted output
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	fmt.Fprintln(w, "ID\tNAME\tLAUNCHER\tPLAYTIME")
	fmt.Fprintln(w, "────\t────────────────────\t──────────\t────────")

	for i, g := range lib.Games {
		fmt.Fprintf(
			w,
			"%d\t%s\t%s\t%dh\n",
			i,
			g.DisplayName,
			g.Launcher,
			g.Playtime / 3600,
		)
	}

	w.Flush()

	return nil
}

func flatGames(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	lib := config.LoadLibrary()
	res := []string{}

	for _, game := range lib.Games {
		res = append(res, game.DisplayName)
	}

	return res, cobra.ShellCompDirectiveNoFileComp
}
