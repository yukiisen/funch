# funch
A small tool to manage locally installed games.

`funch` manages games from multiple sources such as Steam, Lutris, Wine, and native executables through a single CLI interface.

## Features

- Launch games from:
  - Steam
  - Lutris
  - Wine
  - Native executables
  - LOVE games
- Unified game library
- Show status in discord activity
- Shell completion support
- Optional detached/background launching

## Installation

Simply do:
```
go install github.com/yukiisen/funch@v0.1.0-beta
```

## Usage

### Add a game

`funch` will try to guess the proper launcher from the executable name, if none, you have to specify it manually in the config file at (`~/.config/funch/config.toml`)

```bash
funch add "Celeste" ~/games/celeste/Celeste.exe
```

### Run a game

```bash
funch run Celeste
```

### Show game info

```bash
funch info Celeste
```

### List installed games

```bash
funch list
```

### Modify game settings

```bash
funch set RPCEnabled true
```

## Launchers

Supported launchers:

| Launcher | Description                |
| -------- | -------------------------- |
| native   | Run executable directly    |
| wine     | Launch through Wine        |
| steam    | Launch through Steam AppID (still incomplete) |
| lutris   | Launch through Lutris (still incomplete)      |
| love     | Launch LOVE2D games        |

## Shell Completion

### Fish

Generate completions:

```bash
funch completion fish > ~/.config/fish/completions/funch.fish
```

Restart your shell afterwards.

## Daemon Mode

Games can optionally be launched in detached/background mode.

```bash
funch run GameName --daemon
```

This allows helper services such as Discord RPC integration or monitoring to continue running after launch.

## Planned Features

* Steam library syncing
* Lutris game importing
* Playtime statistics
* Tags and filtering
* Game installation helpers (especially for manually installed games)
* Game benshmarking using custom tools
* Extensible launcher system

## License

MIT
