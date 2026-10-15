package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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

	location := locationAreasResponse{}

	err = json.Unmarshal(body, &location)
	if err != nil {
		return locationAreasResponse{}, fmt.Errorf("unmarshal call went haywire: %w", err)
	}

	return location, nil

}
