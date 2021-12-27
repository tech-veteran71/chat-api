package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	ChatAPI  ConfigChatAPI  `json:"chat-api"`
	Database ConfigDatabase `json:"database"`
	Proxy    string         `json:"proxy"`
}

type ConfigChatAPI struct {
	URL     string `json:"url"`
	Token   string `json:"token"`
	Webhook string `json:"webhook"`
}

type ConfigDatabase struct {
	Driver string `json:"driver"`
	DSN    string `json:"dsn"`
}

func ReadConfig(path string) (*Config, error) {
	var config Config

	// Open configuration file.
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Read configuration file.
	err = json.NewDecoder(f).Decode(&config)
	if err != nil {
		return nil, err
	}

	// Set default values.
	if config.Database.Driver == "" {
		config.Database.Driver = "sqlite3"
	}
	if config.Database.DSN == "" {
		config.Database.DSN = "data.sqlite"
	}

	// Configuration read.
	return &config, nil
}
