package main

import (
	"fmt"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	funcs map[string]func(*state, command) error
}

func (c *commands) Run(s *state, cmd command) error {
	cmdFunc, exists := c.funcs[cmd.Name]
	if !exists {
		return fmt.Errorf("unknown command: %v", cmd.Name)
	}
	err := cmdFunc(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *commands) Register(name string, f func(*state, command) error) error {
	if _, exists := c.funcs[name]; exists {
		return fmt.Errorf("command %v already exists", name)
	}
	c.funcs[name] = f
	return nil
}

func NewCommands() commands {
	return commands{
		funcs: make(map[string]func(*state, command) error),
	}
}
