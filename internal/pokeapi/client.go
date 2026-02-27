package pokeapi

import (
	"encoding/json"
	"errors"
	"math/rand/v2"

	//"errors"
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
type PokemonAtLocation struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	IsDefault      bool   `json:"is_default"`
	Order          int    `json:"order"`
	Weight         int    `json:"weight"`
	Abilities      []struct {
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
		Ability  struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
	} `json:"abilities"`
	Forms []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	Species struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
}

func (c *Client) GetPokeLocations(url string) ([]string, [2]string, error) {
	if url == "" {
		url = BASEURL
	}
	var locations = Location{}
	if value, ok := c.cache.Get(url); ok {
		err := json.Unmarshal(value, &locations)
		if err != nil {
			log.Fatal(err)
		}
	} else {
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
		c.cache.Add(url, body)
		err = json.Unmarshal(body, &locations)

		if err != nil {
			log.Fatalf("Error unmarshaling JSON: %s", err.Error())
		}
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

func (c *Client) GetAreaPokemon(location string) (areaPokemon []string, err error) {
	if len(location) == 0 {
		return areaPokemon, errors.New("No location selected")
	}
	var locationPokemonStruct PokemonAtLocation
	url := BASEURL + "/" + location
	if value, ok := c.cache.Get(url); ok {
		err := json.Unmarshal(value, &locationPokemonStruct)
		if err != nil {
			log.Fatal(err)
		}
	} else {
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
		c.cache.Add(url, body)
		err = json.Unmarshal(body, &locationPokemonStruct)

		if err != nil {
			log.Fatalf("Error unmarshaling JSON: %s", err.Error())
		}
	}

	for l := range locationPokemonStruct.PokemonEncounters {
		areaPokemon = append(areaPokemon, locationPokemonStruct.PokemonEncounters[l].Pokemon.Name)
	}

	return areaPokemon, err
}
func (c *Client) TryToCatchPokemon(pokemon string) (caughtPokemon Pokemon, err error) {
	if len(pokemon) == 0 {
		return Pokemon{}, errors.New("No pokemon selected")
	}
	url := POKEMONURL + "/" + pokemon
	var pokemonStruct Pokemon
	if value, ok := c.cache.Get(url); ok {
		err := json.Unmarshal(value, &pokemonStruct)
		if err != nil {
			log.Fatal(err)
		}
	} else {
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
		c.cache.Add(url, body)
		err = json.Unmarshal(body, &pokemonStruct)

		if err != nil {
			log.Fatalf("Error unmarshaling JSON: %s", err.Error())
		}
	}

	caught := tryToCatchPokemon(pokemonStruct.BaseExperience)
	if caught {
		caughtPokemon := pokemonStruct
		return caughtPokemon, nil
	}
	return Pokemon{}, nil
}

func tryToCatchPokemon(exp int) bool {
	roll := rand.IntN(exp)
	if roll <= 30 {
		return true
	}
	return false
}
