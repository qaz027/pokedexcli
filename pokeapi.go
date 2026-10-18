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

const locationURL = "https://pokeapi.co/api/v2/location-area/"
const pokemonURL = "https://pokeapi.co/api/v2/pokemon/"

type Pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
}

type locationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type pokemonRef struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type locationAreasResponse struct {
	Count    int            `json:"count"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
	Results  []locationArea `json:"results"`
}

type locationAreaResponse struct {
	PokemonEncounters []encounter `json:"pokemon_encounters"`
}

type encounter struct {
	Pokemon pokemonRef `json:"pokemon"`
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

func fetchLocationArea(name string) (locationAreaResponse, error) {

	url := locationURL + name

	cached, ok := cache.Get(url)

	if ok {
		location := locationAreaResponse{}

		err := json.Unmarshal(cached, &location)
		if err != nil {
			return locationAreaResponse{}, fmt.Errorf("unmarshal call went haywire: %w", err)
		}
		return location, nil
	}

	res, err := http.Get(url)

	if err != nil {
		return locationAreaResponse{}, fmt.Errorf("error with Get response: %w", err)

	}

	defer res.Body.Close()

	if res.StatusCode > 299 {
		return locationAreaResponse{}, fmt.Errorf("response failed with status code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return locationAreaResponse{}, fmt.Errorf("error with ReadAll: %w", err)

	}
	//save body to the cache
	cache.Add(url, body)

	location := locationAreaResponse{}

	err = json.Unmarshal(body, &location)
	if err != nil {
		return locationAreaResponse{}, fmt.Errorf("unmarshal call went haywire: %w", err)
	}

	return location, nil

}

func fetchPokemon(name string) (Pokemon, error) {

	url := pokemonURL + name

	cached, ok := cache.Get(url)

	if ok {
		pokemon := Pokemon{}

		err := json.Unmarshal(cached, &pokemon)
		if err != nil {
			return Pokemon{}, fmt.Errorf("unmarshal call went haywire: %w", err)
		}
		return pokemon, nil
	}

	res, err := http.Get(url)

	if err != nil {
		return Pokemon{}, fmt.Errorf("error with Get response: %w", err)

	}

	defer res.Body.Close()

	if res.StatusCode > 299 {
		return Pokemon{}, fmt.Errorf("response failed with status code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return Pokemon{}, fmt.Errorf("error with ReadAll: %w", err)

	}
	//save body to the cache
	cache.Add(url, body)

	pokemon := Pokemon{}

	err = json.Unmarshal(body, &pokemon)
	if err != nil {
		return Pokemon{}, fmt.Errorf("unmarshal call went haywire: %w", err)
	}

	return pokemon, nil

}
