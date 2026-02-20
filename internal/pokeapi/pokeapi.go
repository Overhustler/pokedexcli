package pokeapi

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Location struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func GetPokeLocations(url string) ([]string, [2]string, error) {
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area"
	}

	res, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()

	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}

	if err != nil {
		log.Fatal(err)
	}
	var locations = Location{}
	err = json.Unmarshal(body, &locations)

	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %s", err.Error())
	}

	var locationsSlice []string

	for _, l := range locations.Results {
		locationsSlice = append(locationsSlice, l.Name)
	}

	urls := [2]string{*locations.Next, *locations.Previous}

	return locationsSlice, urls, nil

}
