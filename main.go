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
	callback    func(*config, []string) error
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
		"explore": {
			name:        "explore",
			description: "Displays the named location",
			callback:    commandExplore,
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
			runcmd := command.callback(cfg, words[1:])
			if runcmd != nil {
				fmt.Println(runcmd)
			}
			continue
		}
		fmt.Println("Unknown command")

	}

}

func commandExit(cfg *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, args []string) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	fmt.Println()
	for _, cmd := range cfg.commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(cfg *config, args []string) error {
	url := baseURL
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

func commandMapb(cfg *config, args []string) error {
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

func commandExplore(cfg *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("you must provide a location name")
	}

	name := args[0]
	fmt.Printf("Exploring %s...\n", name)
	page, err := fetchLocationArea(name)
	if err != nil {
		return fmt.Errorf("error with fetch: %w", err)
	}
	fmt.Println("Found Pokemon:")
	for _, enc := range page.PokemonEncounters {
		fmt.Printf(" - %s\n", enc.Pokemon.Name)
	}
	return nil
}
