package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlexTLDR/GoPokedex/internal/pokeapi"
)

type commandFunc func(cfg *config) error

type commandWithDescription struct {
	exec        commandFunc
	description string
}

var commands map[string]commandWithDescription
var client pokeapi.Client

func init() {
	client = pokeapi.New()

	commands = map[string]commandWithDescription{
		"map": {
			exec:        commandMap,
			description: "Display the names of 20 location areas in the Pokemon world",
		},
		"mapb": {
			exec:        commandMapBack,
			description: "Go back to the previous list of location areas",
		},
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

func main() {
	cfg := &config{}

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
	locationResp, err := client.ListLocationAreas(cfg.Next)
	if err != nil {
		return err
	}

	if locationResp.Next != nil {
		cfg.Next = *locationResp.Next
	} else {
		cfg.Next = ""
	}

	if locationResp.Previous != nil {
		cfg.Previous = *locationResp.Previous
	} else {
		cfg.Previous = ""
	}

	for _, location := range locationResp.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandMapBack(cfg *config) error {
	if cfg.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	locationResp, err := client.ListLocationAreas(cfg.Previous)
	if err != nil {
		return err
	}

	if locationResp.Next != nil {
		cfg.Next = *locationResp.Next
	} else {
		cfg.Next = ""
	}

	if locationResp.Previous != nil {
		cfg.Previous = *locationResp.Previous
	} else {
		cfg.Previous = ""
	}

	for _, location := range locationResp.Results {
		fmt.Println(location.Name)
	}

	return nil
}
