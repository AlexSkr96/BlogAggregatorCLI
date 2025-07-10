package main

import (
	"fmt"
	"os"

	"github.com/AlexSkr96/BlogAggregatorCLI/internal/config"
)

func main() {
	var state config.State
	var err error
	state.Config, err = config.ReadConfig()
	if err != nil {
		fmt.Printf("ERROR while getting config from file: %v\n", err)
		os.Exit(1)
	}

	commands := config.NewCommands()
	commands.Register("login", config.HandlerLogin)

	if len(os.Args) < 2 {
		fmt.Printf("ERROR: not enough arguments were provided\n")
		os.Exit(1)
	}
	command := config.Command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}
	err = commands.Run(&state, command)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
