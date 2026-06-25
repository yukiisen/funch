package api

import "time"

type Game struct {
	Name           string    `toml:"name"`
	DisplayName    string    `toml:"displayname"`
	Genres         []string  `toml:"genres"`
	Platforms      []string  `toml:"platforms"`
	Stores         []string  `toml:"stores"`
	Website        string    `toml:"website"`
	Icons     	   []string  `toml:"icons"`
	IconIdx		   int       `toml:"iconidx"`
	Launcher       string    `toml:"launcher"`
	Args           []string  `toml:"args"`
	Executable     string    `toml:"executable"`
	Playtime       int64     `toml:"playtime"`
	Installed      time.Time `toml:"installed"`
}


type GameProvider interface {
	GetGame(string) (Game, error)
}
