package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DbUrl    string "json:db_url"
	Username string "json:current_user_name"
}

const configFileName = ".gatorconfig.json"

func ReadConfig() (Config, error) {
	path := "./" + configFileName
	rawJson, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var config Config
	err = json.Unmarshal(rawJson, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
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
