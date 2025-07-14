package main

import (
	"context"
	"fmt"
)

func RegisterHandler(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("expecting a username")
	}
	user, err := s.db.GetUserByUsername(context.Background(), cmd.Args[0])
	if err != nil {
		return err
	}
	if user != nil {
		fmt.Printf("User %s already exists\n", cmd.Args[0])
		return nil
	}
	s.db.CreateUser(cmd.Args[0])
	return nil
}
