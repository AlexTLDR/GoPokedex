package main

import (
	"bufio"
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
		"explore": {
			exec:        commandExplore,
			description: "Explore the Pokemon world",
		},
	}
}

type config struct {
	Next     string
	Previous string
	Args     []string
}

func main() {
	cfg := &config{}
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())

		if input == "exit" {
			break
		}

		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}

		cmdName := parts[0]
		args := parts[1:]

		if cmd, exists := commands[cmdName]; exists {
			err := executeCommand(cmd, cfg, args)
			if err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			fmt.Println("Unknown command. Type 'help' for a list of commands.")
		}
	}
}

func executeCommand(cmd commandWithDescription, cfg *config, args []string) error {
	cfg.Args = args
	return cmd.exec(cfg)
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

func commandExplore(cfg *config) error {
	if len(cfg.Args) == 0 {
		return fmt.Errorf("please provide a location area name or id to explore")
	}

	locationName := cfg.Args[0]
	fmt.Printf("Exploring %s...\n", locationName)

	locationArea, err := client.GetLocationArea(locationName)
	if err != nil {
		return fmt.Errorf("error exploring location area: %w", err)
	}

	fmt.Printf("Pokemon in %s:\n", locationArea.Name)
	for _, pokemon := range locationArea.PokemonEncounters {
		fmt.Println(pokemon.Pokemon.Name)
	}
	return nil
}
