package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Overhustler/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(c *config, areaName ...string) error
}

type config struct {
	pokeApiClient pokeapi.Client
	nextURL       string
	previousURL   string
	pokedex       map[string]pokeapi.Pokemon
}

func buildCommands() map[string]cliCommand {
	commands := map[string]cliCommand{}

	commands["exit"] = cliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    CommandExit,
	}

	commands["help"] = cliCommand{
		name:        "help",
		description: "Show available commands",
		callback: func(c *config, location ...string) error {
			return CommandHelp(c, commands)
		},
	}

	commands["map"] = cliCommand{
		name:        "map",
		description: "display location areas",
		callback:    CommandMap,
	}
	commands["mapb"] = cliCommand{
		name:        "mapb",
		description: "display previous location areas",
		callback:    CommandMapB,
	}
	commands["explore"] = cliCommand{
		name:        "explore",
		description: "display pokemon in the area",
		callback:    CommandExplore,
	}
	commands["catch"] = cliCommand{
		name:        "catch",
		description: "try to catch pokemon",
		callback:    CommandCatch,
	}
	commands["inspect"] = cliCommand{
		name:        "inspect",
		description: "inspect a pokemon you have caught",
		callback:    CommandInspect,
	}
	return commands
}

func repl() {
	scanner := bufio.NewScanner(os.Stdin)
	pokeClient := pokeapi.NewClient(5*time.Second, 5*time.Minute)

	cfg := &config{
		pokeApiClient: pokeClient,
		nextURL:       "",
		previousURL:   "",
		pokedex:       map[string]pokeapi.Pokemon{},
	}

	commands := buildCommands()

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		input := CleanInput(scanner.Text())
		if len(input) == 0 {
			continue
		}
		cmd, ok := commands[input[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		if len(input) > 1 {
			if err := cmd.callback(cfg, input[1]); err != nil {
				fmt.Println(err)
			}
		} else {
			if err := cmd.callback(cfg); err != nil {
				fmt.Println(err)
			}
		}
	}
}

func CleanInput(text string) []string {
	fields := strings.Fields(strings.ToLower(text))
	return fields
}

func CommandExit(c *config, location ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return errors.New("Did not exit correctly")
}

func CommandHelp(c *config, commands map[string]cliCommand, location ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")

	if len(commands) == 0 {
		return errors.New("No valid commands")
	}
	for key, value := range commands {
		fmt.Printf("%s: %s\n", key, value.description)
	}
	return nil
}
func CommandMap(c *config, location ...string) error {
	locations, urls, err := c.pokeApiClient.GetPokeLocations(c.nextURL)
	if err != nil {
		log.Fatal(err)
	}
	for _, l := range locations {
		fmt.Println(l)
	}
	c.nextURL = urls[0]
	c.previousURL = urls[1]
	return nil
}
func CommandMapB(c *config, location ...string) error {
	if c.previousURL == "" {
		println("You are on the first page")
		return nil
	}
	locations, urls, err := c.pokeApiClient.GetPokeLocations(c.previousURL)

	if err != nil {
		log.Fatal(err)
	}
	for _, l := range locations {
		fmt.Println(l)
	}
	c.nextURL = urls[0]
	c.previousURL = urls[1]
	return nil
}

func CommandExplore(c *config, location ...string) error {
	if len(location) == 0 {
		return errors.New("no location provided")
	}
	pokemon, err := c.pokeApiClient.GetAreaPokemon(location[0])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Exploring %s...\n", location[0])
	println("Found Pokemon:")
	for _, l := range pokemon {
		fmt.Printf(" - %s\n", l)
	}
	return nil
}

func CommandCatch(c *config, pokemon ...string) error {
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon[0])
	pokemonStruct, err := c.pokeApiClient.TryToCatchPokemon(pokemon[0])
	if err != nil {
		log.Fatal(err)
	}
	if pokemonStruct.Name == "" {
		fmt.Printf("You did not catch %s\n", pokemon[0])
		return nil
	}
	c.pokedex[pokemon[0]] = pokemonStruct
	fmt.Printf("You caught %s\n", pokemon[0])
	return nil
}
func CommandInspect(c *config, pokemon ...string) error {
	pokemonInfo, ok := c.pokedex[pokemon[0]]
	if !ok {
		fmt.Printf("You have not caught %s\n", pokemon[0])
		return nil
	}
	println(formatField(pokemonInfo.Name))
	println(formatField(strconv.Itoa(pokemonInfo.Height)))
	println(formatField(strconv.Itoa(pokemonInfo.Weight)))
	println("stats:")
	for i := range pokemonInfo.Stats {
		println(formatStats(pokemonInfo.Stats[i].Stat.Name, strconv.Itoa(pokemonInfo.Stats[0].BaseStat)))
	}
	for i := range pokemonInfo.Types {
		println(formatType(pokemonInfo.Types[i].Type.Name))
	}
	return nil
}
func CommandPokedex(c *config, lnput ...string) error {
	if len(c.pokedex) == 0 {
		println("You have not caught any Pokemon")
		return nil
	}

	println("Your Pokedex:")

	for i := range c.pokedex {
		fmt.Println(formatType(c.pokedex[i].Name))
	}
	return nil
}
func formatType(pokemonType string) string {
	nameVal := strings.Split(pokemonType, ":")
	return fmt.Sprintf("\t-%s", nameVal[1])
}
func formatStats(statName string, statValue string) string {
	name := strings.Split(statName, ":")
	val := strings.Split(statValue, ":")

	return fmt.Sprintf("\t-%s: %s", name[1], val[1])
}
func formatField(toSplit string) string {
	stringToFormat := strings.Split(toSplit, ":")
	formattedString := fmt.Sprintf("%s: %s", stringToFormat[0], stringToFormat[1])

	return formattedString
}
