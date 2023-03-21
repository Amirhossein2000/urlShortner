package main

import "os"

type Config struct {
	Server ServerConfig
}

type ServerConfig struct {
	Address string
}

func GetConfig() *Config {
	address := os.Getenv("SERVER_ADDRESS")
	if address == "" {
		address = ":8080"
	}
	return &Config{Server: ServerConfig{Address: address}}
}
