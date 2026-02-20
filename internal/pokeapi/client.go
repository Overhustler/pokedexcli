package pokeapi

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Overhustler/pokedexcli/internal/pokecache"
)

type Client struct {
	cache      *pokecache.Cache
	httpClient http.Client
}

func NewClient(timeout, cacheInterval time.Duration) Client {
	return Client{
		cache:      pokecache.NewCache(cacheInterval),
		httpClient: http.Client{Timeout: timeout},
	}
}

type Location struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func (c *Client) GetPokeLocations(url string) ([]string, [2]string, error) {
	if url == "" {
		url = BASEURL
	}
	var locations = Location{}
	if value, ok := c.cache.Get(url); ok {
		err = json.Unmarshal(value, &locations)
	}
	res, err := c.httpClient.Get(url)

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
	err = json.Unmarshal(body, &locations)

	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %s", err.Error())
	}

	var locationsSlice []string

	for l := range locations.Results {
		locationsSlice = append(locationsSlice, locations.Results[l].Name)
	}
	next, previous := "", ""

	if locations.Next != nil {
		next = *locations.Next
	}
	if locations.Previous != nil {
		previous = *locations.Previous
	}
	urls := [2]string{next, previous}
	return locationsSlice, urls, nil

}
