package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DbUrl    string `json:"db_url"`
	Username string `json:"current_user_name"`
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
	fmt.Printf("%v\n", c)
	rawJson, err := json.Marshal(c)
	if err != nil {
		return err
	}
	os.WriteFile(path, rawJson, 0644)
	return nil
}
