package api

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

type RAWGGame struct {
	Name        string         `json:"name"`
	Slug        string         `json:"slug"`
	Website     string         `json:"website"`

	Platforms []struct {
		Platform struct {
			Name string `json:"name"`
		} `json:"platform"`
	} `json:"platforms"`

	Genres []struct {
		Name string `json:"name"`
	} `json:"genres"`

	Stores []struct {
		Store struct {
			Name string `json:"name"`
		} `json:"store"`
	} `json:"stores"`
}

type RAWGResponse struct {
	Results []RAWGGame `json:"results"`
}

const base string = "https://api.rawg.io/api/games"

type RAWG struct {
	Key string
}

func (self *RAWG) GetGame(name string) (Game, error) {
	// prepare a request 
	params := url.Values {}
	params.Add("key", self.Key)
	params.Add("search", name)
	params.Add("page_size", "1")
	requrl := base + "?" + params.Encode()

	var search_res RAWGResponse
	err := requestJson(requrl, &search_res)
	if err != nil { return Game{}, err }

	// check if the game matches
	if len(search_res.Results) == 0 { return Game{}, errors.New("No Results") }
	rgame := search_res.Results[0]

	if rgame.Name != name && rgame.Slug != name {
		input := ask("Did you mean " + rgame.Name + "? [Y/n] (yes)")
		if input == "n" || input == "no" { return Game{}, errors.New("Game Not Found") }
	}

	// game matches, get extra info
	params = url.Values {}
	params.Add("key", self.Key)
	requrl = base + "/" + rgame.Slug + "?" + params.Encode()

	var game_res RAWGGame
	err = requestJson(requrl, &game_res)
	if err != nil { return Game{}, err }

	return makeGame(game_res), nil
}

func ask(prompt string) string {
	fmt.Print(prompt)

	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	return strings.ToLower(strings.TrimSpace(text))
}

func requestJson(site string, target any) error {
	res, err := http.Get(site)
	if err != nil { return err }
	if res.StatusCode != http.StatusOK { return errors.New(res.Status)}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(target)
	if err != nil { return err }

	return nil
}

func makeGame(base RAWGGame) Game {
	game := Game {
		Name: base.Slug,
		DisplayName: base.Name,
		Website: base.Website,
	}

	for _, genre := range base.Genres {
		game.Genres = append(game.Genres, genre.Name)
	}

	for _, store := range base.Stores {
		game.Stores = append(game.Stores, store.Store.Name)
	}

	for _, platform := range base.Platforms {
		game.Platforms = append(game.Platforms, platform.Platform.Name)
	}

	return game
}
