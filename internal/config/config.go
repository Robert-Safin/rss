package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	err := write(*c)
	if err != nil {
		return err
	}
	return nil
}

func Read() (Config, error) {
	var config Config

	pathToConfig, err := getConfigFilePath()
	if err != nil {
		return config, err
	}

	if _, err := os.Stat(pathToConfig); os.IsNotExist(err) {
		fmt.Println("Config file not found, creating default config")

		defaultConfig := Config{
			DbURL:           "postgres://example",
			CurrentUserName: "postgres://rob:@localhost:5432/gator?sslmode=disable",
		}

		if err := write(defaultConfig); err != nil {
			return config, fmt.Errorf("Couldn't create config file: %w", err)
		}

		return defaultConfig, nil
	}

	bytes, err := os.ReadFile(pathToConfig)
	if err != nil {
		return config, fmt.Errorf("Path error: %w", err)
	}

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return config, fmt.Errorf("Unmarshal error: %w", err)
	}

	return config, nil
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("OS error : %w", err)
	}
	pathToConfig := home + "/.gatorconfig.json"

	return pathToConfig, nil
}

func write(config Config) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal error : %w", err)
	}

	err = os.WriteFile(path, bytes, 0644)
	if err != nil {
		return fmt.Errorf("OS write error : %w", err)
	}
	return nil
}
