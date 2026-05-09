package api

import "time"

type Game struct {
	Name           string    `toml:"name"`
	DisplayName    string    `toml:"displayname"`
	Genres         []string  `toml:"genres"`
	Platforms      []string  `toml:"platforms"`
	Stores         []string  `toml:"stores"`
	Website        string    `toml:"website"`
	Launcher       string    `toml:"launcher"`
	Args           []string  `toml:"args"`
	Executable     string    `toml:"executable"`
	Playtime       int       `toml:"playtime"`
	Installed      time.Time `toml:"installed"`
}


type GameProvider interface {
	GetGame(string) (Game, error)
}
