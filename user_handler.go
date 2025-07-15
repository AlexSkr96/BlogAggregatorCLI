package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/AlexSkr96/BlogAggregatorCLI/internal/database"
	"github.com/google/uuid"
)

func HandlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("expecting a username")
	}
	_, err := s.db.GetUserByUsername(context.Background(), cmd.Args[0])
	if err != nil {
		return err
	}
	err = s.config.SetUser(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("error while setting username: %v", err)
	}
	fmt.Printf("user is set")
	return nil
}

func HandlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("expecting a username")
	}
	_, err := s.db.GetUserByUsername(context.Background(), cmd.Args[0])
	if err != nil {
		// If error is not "no rows found", return the error
		if err != sql.ErrNoRows {
			return err
		}
		// User not found, continue with registration
	} else {
		// User found, user already exists
		return fmt.Errorf("User %s already exists\n", cmd.Args[0])
	}
	id := uuid.New()
	createUserParams := database.CreateUserParams{
		ID:        id,
		Name:      cmd.Args[0],
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	s.db.CreateUser(context.Background(), createUserParams)
	fmt.Printf("created user %v with ID %v", cmd.Args[0], id)
	return nil
}
