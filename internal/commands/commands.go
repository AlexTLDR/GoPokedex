package commands

import (
	"fmt"
	"github.com/AlexTLDR/GoPokedex/internal/pokeapi"
	"math/rand/v2"
	"os"
)

type CommandFunc func(cfg *Config) error

type Command struct {
	Exec        CommandFunc
	Description string
}

type Config struct {
	Next          string
	Previous      string
	Args          []string
	CaughtPokemon map[string]pokeapi.Pokemon
}

var Commands map[string]Command
var client pokeapi.Client

func Initialize() {
	client = pokeapi.New()

	Commands = map[string]Command{
		"map": {
			Exec:        CommandMap,
			Description: "Display the names of 20 location areas in the Pokemon world",
		},
		"mapb": {
			Exec:        CommandMapBack,
			Description: "Go back to the previous list of location areas",
		},
		"help": {
			Exec:        CommandHelp,
			Description: "Displays a help message",
		},
		"exit": {
			Exec:        CommandExit,
			Description: "Exit the Pokedex",
		},
		"explore": {
			Exec:        CommandExplore,
			Description: "Explore the Pokemon world",
		},
		"catch": {
			Exec:        CommandCatch,
			Description: "Attempt to catch a Pokemon",
		},
		"inspect": {
			Exec:        CommandInspect,
			Description: "Inspect a caught Pokemon",
		},
	}
}

func ExecuteCommand(cmdName string, cfg *Config, args []string) error {
	if cmd, exists := Commands[cmdName]; exists {
		cfg.Args = args
		return cmd.Exec(cfg)
	}
	return fmt.Errorf("unknown command: %s", cmdName)
}

func CommandExit(cfg *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func CommandHelp(cfg *Config) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for name, cmd := range Commands {
		fmt.Printf("%s: %s\n", name, cmd.Description)
	}
	return nil
}

func CommandMap(cfg *Config) error {
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

func CommandMapBack(cfg *Config) error {
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

func CommandExplore(cfg *Config) error {
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

func CommandCatch(cfg *Config) error {
	if len(cfg.Args) == 0 {
		return fmt.Errorf("please provide a pokemon name")
	}

	pokemonName := cfg.Args[0]

	if _, ok := cfg.CaughtPokemon[pokemonName]; ok {
		return fmt.Errorf("you've already caught %s", pokemonName)
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	pokemon, err := client.GetPokemon(pokemonName)
	if err != nil {
		return err
	}

	catchRate := 0.5 - float64(pokemon.BaseExperience)/1000.0
	if catchRate < 0.1 {
		catchRate = 0.1
	}

	r := rand.Float64()
	if r <= catchRate {
		cfg.CaughtPokemon[pokemonName] = pokemon
		fmt.Printf("%s was caught!\n", pokemonName)
		return nil
	}
	fmt.Printf("%s escaped!\n", pokemonName)
	return nil
}

func CommandInspect(cfg *Config) error {
	if len(cfg.Args) == 0 {
		return fmt.Errorf("please provide a pokemon name")
	}
	pokemonName := cfg.Args[0]
	pokemon, ok := cfg.CaughtPokemon[pokemonName]
	if !ok {
		return fmt.Errorf("you haven't caught %s yet", pokemonName)
	}
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)

	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Println("Types:")
	for _, typeInfo := range pokemon.Types {
		fmt.Printf("  -%s\n", typeInfo.Type.Name)
	}
	return nil
}
