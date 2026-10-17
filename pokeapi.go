package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pokedexcli/internal/pokecache"
	"time"
)

var cache = pokecache.NewCache(5 * time.Second)

type locationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type locationAreasResponse struct {
	Count    int            `json:"count"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
	Results  []locationArea `json:"results"`
}

func fetchLocationAreas(url string) (locationAreasResponse, error) {

	cached, ok := cache.Get(url)

	if ok {
		location := locationAreasResponse{}

		err := json.Unmarshal(cached, &location)
		if err != nil {
			return locationAreasResponse{}, fmt.Errorf("unmarshal call went haywire: %w", err)
		}
		return location, nil
	}

	res, err := http.Get(url)

	if err != nil {
		return locationAreasResponse{}, fmt.Errorf("error with Get response: %w", err)

	}

	defer res.Body.Close()

	if res.StatusCode > 299 {
		return locationAreasResponse{}, fmt.Errorf("response failed with status code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return locationAreasResponse{}, fmt.Errorf("error with ReadAll: %w", err)

	}
	//save body to the cache
	cache.Add(url, body)

	location := locationAreasResponse{}

	err = json.Unmarshal(body, &location)
	if err != nil {
		return locationAreasResponse{}, fmt.Errorf("unmarshal call went haywire: %w", err)
	}

	return location, nil

}
