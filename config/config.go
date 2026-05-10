package config

import (
	"github.com/yukiisen/funch/api"
	"os"
	"path"

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
	FuzzyFinderCmd []string
}

const configDir = "."

func InitLogging (name string) (*os.File, error) {
	err := os.MkdirAll(path.Join(configDir, "logs"), 0755)
	if err != nil { return nil, err }
	
	return os.Create(path.Join(configDir, "logs", name))
}

func createLibrary() Library {
	config := Config {
		RPCEnabled: true,
		RPCShowPlaytime: true,
		Daemonize: false,
		FuzzyFinderCmd: []string{ "bash", "-c", "false" },
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

	_, e := os.Stat(path.Join(configDir , "config.toml"))
	if e != nil {
		return createLibrary()
	}

	_, err := toml.DecodeFile(path.Join(configDir , "config.toml"), &lib)
	if err != nil { panic(err) }

	return lib
}

func UpdateLibrary(lib Library) error {
	file, err := os.Create(path.Join(configDir , "config.toml"));
	if err != nil { return err }
	defer file.Close()

	err = toml.NewEncoder(file).Encode(lib)
	if err != nil { return err }

	return nil
}
