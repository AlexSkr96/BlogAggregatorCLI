package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DbUrl    string "json:db_url"
	Username string "json:current_user_name"
}

<<<<<<< HEAD
type State struct {
=======
type state struct {
>>>>>>> 579614f9ff3598c2c46c2dc6c7d097ae97797cb4
	Config *Config
}

type command struct {
	Name string
	Args []string
}

type commands struct {
	funcs map[string]func(*state, command) error
}

const configFileName = ".gatorconfig.json"

func ReadConfig() (*Config, error) {
	path := "./" + configFileName
	rawJson, err := os.ReadFile(path)
	if err != nil {
		return &Config{}, err
	}
	var config Config
	err = json.Unmarshal(rawJson, &config)
	if err != nil {
		return &Config{}, err
	}
	return &config, nil
}

func (c *Config) SetUser(username string) error {
	path := "./" + configFileName
	c.Username = username
	rawJson, err := json.Marshal(c)
	if err != nil {
		return err
	}
	os.WriteFile(path, rawJson, 0644)
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("expecting a username")
	}
	config, err := ReadConfig()
	if err != nil {
		return fmt.Errorf("error while getting username: %v", err)
	}
	err = config.SetUser(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("error while setting username: %v", err)
	}
	fmt.Printf("user is set")
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	cmdFunc, exists := c.funcs[cmd.Name]
	if exists {
		if !exists {
			return fmt.Errorf("unknown command: %v", cmd)
		}
	}
<<<<<<< HEAD
	err := cmdFunc(s, cmd)
=======
	err := cmdFunc(s.Config)
>>>>>>> 579614f9ff3598c2c46c2dc6c7d097ae97797cb4
	if err != nil {
		return fmt.Errorf("")
	}
	return nil
}

<<<<<<< HEAD
func (c *commands) register(name string, f func(*state, command) error) error {
	if _, exists := c.funcs[name]; exists {
		return fmt.Errorf("command %v already exists", name)
	}
	c.funcs[name] = f
	return nil
}
=======
func (c *commands) register(name string, f func(*state, command) error)
>>>>>>> 579614f9ff3598c2c46c2dc6c7d097ae97797cb4
