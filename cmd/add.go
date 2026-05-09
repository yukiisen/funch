package cmd

import (
	"errors"
	"fmt"
	"funch/api"
	"funch/config"
	"os"
	"path/filepath"
	"time"
	"strings"
	"github.com/spf13/cobra"
)

var AddCmd = &cobra.Command {
	Use: "add <name> <executable>",
	Args: cobra.ExactArgs(2),
	RunE: add,
}

func add(cmd *cobra.Command, args []string) error {
	name, exec := args[0], args[1]

	if strings.HasPrefix(args[1], "~") {
		home, _ := os.UserHomeDir()
		exec = filepath.Join(home, exec[1:])
	}

	exec, err := filepath.Abs(exec)
	if err != nil { return err }

	lib := config.LoadLibrary()


	stat, err := os.Stat(exec)
	if err != nil { return err }
	if stat.IsDir() { return errors.New("Executable must be a file") }

	rawg := api.RAWG { Key: lib.Config.APIKey }

	fmt.Println("Fetching Game metadata...")
	game, err := rawg.GetGame(name)
	if err != nil { return err }

	game.Args = []string{}
	game.Executable = exec
	game.Installed = time.Now()
	game.Playtime = 0
	game.Launcher = getLauncher(exec)

	lib.Games = append(lib.Games, game)

	// I won't return the err directly becuase it might be confusing later
	err = config.UpdateLibrary(lib);
	if err != nil { return err }

	return nil
}

func getLauncher(file string) string {
	ext := filepath.Ext(file)

	switch ext {
	case ".exe":
		return "wine"
	case ".love":
		return "love"
	case ".AppImage", "":
		return "native"
	case ".sh":
		return "sh"
	}

	return "none"
}
