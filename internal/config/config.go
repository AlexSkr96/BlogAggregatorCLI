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

type State struct {
	Config *Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	funcs map[string]func(*State, Command) error
}

const configFileName = ".gatorconfig.json"

func ReadConfig() (*Config, error) {
	path, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	path += fmt.Sprintf("/%v", configFileName)
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
	path, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	path += fmt.Sprintf("/%v", configFileName)
	c.Username = username
	rawJson, err := json.Marshal(c)
	if err != nil {
		return err
	}
	os.WriteFile(path, rawJson, 0644)
	return nil
}

func HandlerLogin(s *State, cmd Command) error {
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

func (c *Commands) Run(s *State, cmd Command) error {
	cmdFunc, exists := c.funcs[cmd.Name]
	if exists {
		if !exists {
			return fmt.Errorf("unknown command: %v", cmd)
		}
	}
	err := cmdFunc(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *Commands) Register(name string, f func(*State, Command) error) error {
	if _, exists := c.funcs[name]; exists {
		return fmt.Errorf("command %v already exists", name)
	}
	c.funcs[name] = f
	return nil
}

func NewCommands() Commands {
	return Commands{
		funcs: make(map[string]func(*State, Command) error),
	}
}
