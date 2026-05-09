package config

import (
	"funch/api"
	"os"
	"github.com/BurntSushi/toml"
)

type Library struct {
	Games []api.Game `toml:"games"`
	Config Config
}

// app configuration (for now)
type Config struct {
	RPCEnabled bool
	RPCShowPlaytime bool
	Daemonize bool
	APIKey string
}

const configPath = "./config.toml"

func createLibrary() Library {
	config := Config {
		RPCEnabled: true,
		RPCShowPlaytime: true,
		Daemonize: false,
	}

	lib := Library {
		Games: []api.Game{},
		Config: config,
	}

	UpdateLibrary(lib)

	return lib
}

func LoadLibrary() Library {
	var lib Library

	_, e := os.Stat(configPath)
	if e != nil {
		return createLibrary()
	}

	_, err := toml.DecodeFile(configPath, &lib)
	if err != nil { panic(err) }

	return lib
}

func UpdateLibrary(lib Library) error {
	file, err := os.Create(configPath);
	if err != nil { return err }
	defer file.Close()

	err = toml.NewEncoder(file).Encode(lib)
	if err != nil { return err }

	return nil
}
