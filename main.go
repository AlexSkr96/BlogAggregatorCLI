package main

import (
	"context"
	"fmt"
	"os"

	"database/sql"

	_ "github.com/lib/pq"

	"github.com/AlexSkr96/BlogAggregatorCLI/internal/config"
	"github.com/AlexSkr96/BlogAggregatorCLI/internal/database"
)

type state struct {
	config *config.Config
	db     *database.Queries
}

func main() {
	var state state
	var err error
	state.config, err = config.ReadConfig()
	if err != nil {
		fmt.Printf("ERROR while getting config from file: %v\n", err)
		os.Exit(1)
	}
	dbURL := state.config.DbUrl
	db, err := sql.Open("postgres", dbURL)
	state.db = database.New(db)

	commands := NewCommands()
	commands.Register("login", HandlerLogin)
	commands.Register("register", HandlerRegister)
	commands.Register("reset", HandlerReset)
	commands.Register("users", HandlerGetUsers)
	commands.Register("agg", HandlerAggregate)
	commands.Register("addfeed", middlewareLoggedIn(HandlerNewFeed))
	commands.Register("feeds", HandlerFeeds)
	commands.Register("follow", middlewareLoggedIn(HandlerNewFeedFollow))
	commands.Register("following", middlewareLoggedIn(HandlerGetFollowsForUser))
	commands.Register("unfollow", middlewareLoggedIn(HandlerDeleteFeedFollow))

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

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUserByUsername(context.Background(), s.config.Username)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}
