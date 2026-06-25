package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"time"
	"strings"
	"os"
	"syscall"

	"github.com/yukiisen/funch/api"
	"github.com/yukiisen/funch/config"
	"github.com/yukiisen/funch/discord"
	"github.com/spf13/cobra"
)

type Game = api.Game

var RunCmd = &cobra.Command {
	Use: "run <game>",
	Short: "Launch a game",
	Long:  "Launches a game using its detected or assigned launcher (Steam, Lutris, Wine, or native).",
	Args: cobra.MaximumNArgs(1),
	ValidArgsFunction: flatGames,
	RunE: run,
}

func InitRunCmd() {
	RunCmd.Flags().Bool("no-rpc", false, "Force disable discord RPC")
	RunCmd.Flags().Bool("attached", false, "Don't run in background")
	RunCmd.Flags().BoolP("daemon", "d", false, "Run in background")
}

const RPCAppID string = "1355243267494514903"

func run(cmd *cobra.Command, args []string) error {
	lib := config.LoadLibrary()

	daemon, err := cmd.Flags().GetBool("daemon")  
	if err != nil { return err }

	attached, err := cmd.Flags().GetBool("attached")  
	if err != nil { return err }

	if (lib.Config.Daemonize || daemon) && len(args) != 0 && !attached {
		args := append([]string{}, os.Args...)
		args = append(args, "--attached")
		command := exec.Command(os.Args[0], args[1:]...)

		// detach from terminal
		// command.Stdout = nil
		// command.Stderr = nil
		command.Stdin = nil

		command.SysProcAttr = &syscall.SysProcAttr{
			Setsid: true,
		}

		if err := command.Start(); err != nil {
			return err
		}

		fmt.Println("Game launched with PID:", command.Process.Pid)
		return nil
	}

	var game *Game

	if (len(args) == 0)  {
		game = getGameInteractive(lib.Games, lib.Config.FuzzyFinderCmd)
		if game == nil {
			return errors.New("No Game Selected")
		}
	} else {
		game, err = getGame(args[0], lib.Games)
	}

	if err != nil { return err }

	command, err := buildLaunchCommand(game)
	if err != nil { return err }

	// setup logging
	file, err := config.InitLogging(game.Name + time.Now().Format("2006-01-02@15:04:05"))
	if err != nil { return err }
	defer file.Close()

	command.Stdout = file
	command.Stderr = file

	// start game
	if err := command.Start(); err != nil {
		return err
	}

	// discord RPC
	no_rpc, err := cmd.Flags().GetBool("no-rpc")  
	if err != nil { return err }

	client := discord.Client {}

	if lib.Config.RPCEnabled && !no_rpc {
		if err := client.Connect(RPCAppID); err != nil {
			log.Println("RPC Error: ", err)
		} else {
			defer client.Close()
			if err := startRPC(&client, game, time.Now(), lib.Config.RPCShowPlaytime); err != nil {
				log.Println("RPC Error: ", err)
			}
		}
	}
	

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				game.Playtime += 30
				config.UpdateLibrary(lib)
			case <-done:
				return
			}
		}
	} ()

	/// wait until we're done
	if err := command.Wait(); err != nil {
		return err
	}

	done <- true

	return nil
}

func getGameInteractive(games []Game, command []string) *Game {
	var buf bytes.Buffer

	for _, game := range games {
		fmt.Fprintf(&buf, "[%s] %s\n", game.Launcher, game.DisplayName)
	}

	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdin = &buf

	out, err := cmd.Output()
	if err != nil { return nil }

	name := strings.TrimSpace(strings.Split(string(out), "]")[1])

	for i := range games {
		if name == games[i].DisplayName {
			return &games[i]
		}
	}

	return nil
}

func startRPC(client *discord.Client, game *Game, startTime time.Time, showTime bool) error {
	var state string

	if showTime {
		state = "for " + strconv.FormatInt(game.Playtime / 3600, 10) + "h"
	} else {
		state = ""
	}

	return client.SetActivity(&discord.Activity{ 
		State:      state,
		Details:    "Playing " + game.DisplayName,

		Assets: &discord.Assets{
			LargeImage: game.Icons[game.IconIdx],
			LargeText:  game.DisplayName,

			SmallImage: game.Launcher,
			SmallText:  "Via " + game.Launcher,
		},

		Timestamps: &discord.Timestamps{
			Start: startTime.Unix(),
		},

		Buttons: []*discord.Button{
			{
				Label: "Website",
				URL:   game.Website,
			},
		},
	})
}

func buildLaunchCommand(game *Game) (*exec.Cmd, error) {
	switch game.Launcher {
	case "wine", "sh", "love":
		args := append([]string{"--net=none", game.Executable}, game.Args...)
		return exec.Command("firejail", args...), nil

	case "lutris":
		args := append([]string{"lutris:game/" + game.Executable}, game.Args...)
		return exec.Command("lutris", args...), nil

	case "steam":
		args := append([]string{"steam://rungameid/" + game.Executable}, game.Args...)
		return exec.Command("steam", args...), nil

	case "native":
		args := append([]string{game.Executable}, game.Args...)
		return exec.Command(args[0], args...), nil

	default:
		return nil, errors.New("Could not run this game with any known launcher")
	}
}
