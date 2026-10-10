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
	callback    func() error
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
	}
	return SupportedCommands
}

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	SupportedCommands := getCommands()
	//REPL loop
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
		command, ok := SupportedCommands[words[0]]
		if ok {
			runcmd := command.callback()
			if runcmd != nil {
				fmt.Println(runcmd)
			}
			continue
		}
		fmt.Println("Unknown command")

	}

}

func originalmain() {

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		//get input as string
		scanner.Scan()
		input := scanner.Text()
		words := cleanInput(input) // clean text

		//input := scanner.Scan().Text() // can I combine methods? // no

		//how do I capture the first word input[0]?
		if len(words) == 0 {
			continue
		}
		fmt.Println("Your command was:", words[0])

	}

}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	SupportedCommands := getCommands()
	fmt.Println("Welcome to the Pokedex!\nUsage:\n")
	for _, cmd := range SupportedCommands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}
