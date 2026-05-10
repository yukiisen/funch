package cmd

import (
	"errors"
	"fmt"
	"funch/config"
	"os"
	"text/tabwriter"
	"strconv"
	"strings"
	"github.com/spf13/cobra"
)

var InfoCmd = &cobra.Command {
	Use: "info <game>",
	Short: "Show game details",
	Long:  "Displays detailed information about a specific game, including its launcher type, executable, and configuration data.",
	Args: cobra.ExactArgs(1),
	ValidArgsFunction: flatGames,
	RunE: info,
}

func info(cmd *cobra.Command, args []string) error {
	name := args[0]
	lib := config.LoadLibrary()

	g, err := getGame(name, lib.Games)
	if err != nil { return err }

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	fmt.Printf("%s\n", g.DisplayName)
	fmt.Println("────────")

	fmt.Fprintf(w, "Launcher\t%s\n", g.Launcher)
	fmt.Fprintf(w, "Executable\t%s\n", g.Executable)
	fmt.Fprintf(w, "Installed\t%s\n", g.Installed.Format("2006-01-02"))
	fmt.Fprintf(w, "Playtime\t%dh\n", g.Playtime / 3600)

	fmt.Fprintf(w, "Genres\t%s\n", strings.Join(g.Genres, ", "))
	fmt.Fprintf(w, "Platforms\t%s\n", strings.Join(g.Platforms, ", "))
	fmt.Fprintf(w, "Stores\t%s\n", strings.Join(g.Stores, ", "))

	fmt.Fprintf(w, "Website\t%s\n", g.Website)

	w.Flush()

	return nil
}


// matches game against name if none then display name if none then id
func getGame(name string, games []Game) (*Game, error) {
	for i := range games {
		if name == games[i].Name || name == games[i].DisplayName {
			return &games[i], nil
		}
	}

	i, err := strconv.Atoi(name)

	if err == nil && i >= 0 && i < len(games) {
		return &games[i], nil
	}

	return nil, errors.New("Game not found")
}
