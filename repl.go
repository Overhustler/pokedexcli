package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(c *config) error
}

type config struct {
	nextURL     string
	previousURL string
}

func buildCommands(c *config) map[string]cliCommand {
	commands := map[string]cliCommand{}

	commands["exit"] = cliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    CommandExit,
	}

	commands["help"] = cliCommand{
		name:        "help",
		description: "Show available commands",
		callback: func(c *config) error {
			return CommandHelp(c, commands)
		},
	}

	commands["map"] = cliCommand{
		name:        "map",
		description: "display location areas",
		callback:    CommandMap,
	}

	return commands
}

func repl() {
	scanner := bufio.NewScanner(os.Stdin)
	cfg := &config{}
	commands := buildCommands(cfg)

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

		if err := cmd.callback(cfg); err != nil {
			fmt.Println(err)
		}
	}
}

func CleanInput(text string) []string {
	fields := strings.Fields(strings.ToLower(text))
	return fields
}

func CommandExit(c *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func CommandHelp(c *config, commands map[string]cliCommand) error {
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
func CommandMap(c *config) error {
	return nil
}
