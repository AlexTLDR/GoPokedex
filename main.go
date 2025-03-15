package main

import (
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

type commandFunc func(cfg *config) error

type commandWithDescription struct {
	exec        commandFunc
	description string
}

var commands map[string]commandWithDescription

func init() {
	commands = map[string]commandWithDescription{
		"map": {
			exec:        commandMap,
			description: "Display the names of 20 location areas in the Pokemon world",
		},
		// "mapb": {
		//     exec:        commandMapBack,
		//     description: "Go back to the previous list of location areas",
		// },
		"help": {
			exec:        commandHelp,
			description: "Displays a help message",
		},
		"exit": {
			exec:        commandExit,
			description: "Exit the Pokedex",
		},
	}
}

type config struct {
	Next     string
	Previous string
}

type LocationAreaResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func main() {
	cfg := &config{}

	// PokeFans CLI loop
	for {
		var input string
		fmt.Print("Pokedex > ")
		fmt.Scanln(&input)

		if cmd, exists := commands[input]; exists {
			err := cmd.exec(cfg)
			if err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			fmt.Println("Unknown command. Type 'help' for a list of commands.")
		}
	}
}

func cleanInput(text string) []string {
	text = strings.TrimSpace(text)
	if text == "" {
		return []string{}
	}
	parts := strings.Fields(text)
	var cleaned []string
	for _, part := range parts {
		cleaned = append(cleaned, strings.ToLower(part))
	}
	return cleaned
}

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for name, cmd := range commands {
		fmt.Printf("%s: %s\n", name, cmd.description)
	}
	return nil
}

func commandMap(cfg *config) error {
	return nil
}
