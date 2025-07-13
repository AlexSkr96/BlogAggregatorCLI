package main

import (
	"fmt"
)

func HandlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("expecting a username")
	}
	// config, err := config.ReadConfig()
	// if err != nil {
	// 	return fmt.Errorf("error while getting username: %v", err)
	// }
	err := s.config.SetUser(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("error while setting username: %v", err)
	}
	fmt.Printf("user is set")
	return nil
}
