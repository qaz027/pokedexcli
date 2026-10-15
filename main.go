package main

import (
	"bufio"
	"fmt"
	"os"
)

// create registry of commands
type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	commands map[string]cliCommand
	previous *string
	next     *string
}

// map of supported commands
func getCommands() map[string]cliCommand {
	var SupportedCommands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the list of the next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the list of the previous 20 locations",
			callback:    commandMapb,
		},
	}
	return SupportedCommands
}

func main() {

	cfg := config{
		commands: getCommands(),
	}

	replLoop(&cfg)

}

func replLoop(cfg *config) {

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		//get input as string
		scanner.Scan()
		input := scanner.Text()
		words := cleanInput(input) // clean text
		if len(words) == 0 {
			continue
		}

		//look up the callback function in the registry
		command, ok := cfg.commands[words[0]]
		if ok {
			runcmd := command.callback(cfg)
			if runcmd != nil {
				fmt.Println(runcmd)
			}
			continue
		}
		fmt.Println("Unknown command")

	}

}

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	fmt.Println()
	for _, cmd := range cfg.commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(cfg *config) error {
	url := "https://pokeapi.co/api/v2/location-area/"
	if cfg.next != nil {
		url = *cfg.next
	}

	page, err := fetchLocationAreas(url)
	if err != nil {
		return fmt.Errorf("error with fetch: %w", err)
	}
	cfg.next = page.Next
	cfg.previous = page.Previous
	for _, area := range page.Results {
		fmt.Println(area.Name)
	}

	return nil
}

func commandMapb(cfg *config) error {
	if cfg.previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	page, err := fetchLocationAreas(*cfg.previous)
	if err != nil {
		return fmt.Errorf("error with fetch: %w", err)
	}
	cfg.next = page.Next
	cfg.previous = page.Previous
	for _, area := range page.Results {
		fmt.Println(area.Name)
	}

	return nil
}
