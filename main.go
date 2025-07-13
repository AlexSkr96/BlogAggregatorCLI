package main

import (
	"fmt"
	"os"

	// "github.com/google/uuid"
	_ "github.com/lib/pq"

	// "database/sql"
	"github.com/AlexSkr96/BlogAggregatorCLI/internal/config"
)

type state struct {
	config *config.Config
}

func main() {
	var state state
	var err error
	state.config, err = config.ReadConfig()
	fmt.Printf("%+v\n", state.config)
	if err != nil {
		fmt.Printf("ERROR while getting config from file: %v\n", err)
		os.Exit(1)
	}
	// dbURL :=
	// db, err := sql.Open("postgres", dbURL)

	commands := NewCommands()
	commands.Register("login", HandlerLogin)

	if len(os.Args) < 2 {
		fmt.Printf("ERROR: not enough arguments were provided\n")
		os.Exit(1)
	}
	command := command{
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
