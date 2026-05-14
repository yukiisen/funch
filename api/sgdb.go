package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type SGDB struct {
	Key string
}

type SGDBSearchResult struct {
	Success bool `json:"success"`
	Data    []struct {
		ID int `json:"id"`
	} `json:"data"`
}

type SGDBGridResult struct {
	Success bool `json:"success"`
	Data    []struct {
		URL string `json:"url"`
	} `json:"data"`
}

func (self *SGDB) GetGameIcons(name string) ([]string, error) {
	url := "https://www.steamgriddb.com/api/v2/search/autocomplete/" + name

	search := SGDBSearchResult {}

	err := requestJsonAuth(url, self.Key, &search);
	if err != nil { return nil, err }

	if len(search.Data) == 0 { return nil, errors.New("fuck no game") }

	id := search.Data[0].ID
	
	url = fmt.Sprintf("https://www.steamgriddb.com/api/v2/icons/game/%d", id)
	grids := SGDBGridResult {}

	err = requestJsonAuth(url, self.Key, &grids);
	if err != nil { return nil, err }

	if len(grids.Data) == 0 { return nil, errors.New("fuck no game") }

	ret := []string{}

	for _, d := range grids.Data {
		ret = append(ret, d.URL)
	}

	return ret, nil
}

func requestJsonAuth(site string, key string, target any) error {
	req, err := http.NewRequest("GET", site, nil)
	if err != nil { return err }

	req.Header.Set("Authorization", "Bearer " + key)

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil { return err }
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New(res.Status)
	}

	return json.NewDecoder(res.Body).Decode(target)
}
