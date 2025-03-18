package main

import (
	"bufio"
	"fmt"
	"github.com/AlexTLDR/GoPokedex/internal/pokeapi"
	"os"
	"strings"

	"github.com/AlexTLDR/GoPokedex/internal/commands"
)

func main() {
	commands.Initialize()

	cfg := &commands.Config{
		CaughtPokemon: make(map[string]pokeapi.Pokemon),
	}
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

		err := commands.ExecuteCommand(cmdName, cfg, args)
		if err != nil {
			fmt.Println("Error:", err)
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
