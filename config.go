package main

import (
	"log"

	"github.com/BurntSushi/toml"
)

// Config the crawler
type Config struct {
	URL               string `toml:"url"`
	DeckPath          string `toml:"deck_path"`
	DeckDataPath      string `toml:"deck_data_path"`
	SideBoardPath     string `toml:"side_board_path"`
	SideBoardDataPath string `toml:"side_board_data_path"`
	OtherIndexPath    string `toml:"other_indexes_path"`
}

func (c *Config) loadConfigFile(filePath string) (map[string]Config, error) {
	var configs map[string]Config
	if _, err := toml.DecodeFile(filePath, &configs); err != nil {
		log.Printf("Error reading TOML file: %s", err)
		return nil, err
	}
	return configs, nil
}

func (c *Config) loadConfig(configName string) (Config, error) {
	configs, err := c.loadConfigFile("config.toml")
	if err != nil {
		return Config{}, err
	}

	return configs[configName], nil
}
