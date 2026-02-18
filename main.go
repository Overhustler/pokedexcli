package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex >")
		scanner.Scan()
		if scanner.Err() != nil {
			log.Fatal(scanner.Err())
		}
		userInput := CleanInput(scanner.Text())
		if len(userInput) == 0 {
			fmt.Println("Using fmt.Println:", errors.New("No input found"))
		} else {
			fmt.Printf("Your command was: %s\n", userInput[0])
		}
	}
}
