package main

import (
	"os"
	"urlShortner/internal/consts"
)

type Config struct {
	Server ServerConfig
}

type ServerConfig struct {
	Address string
}

func GetConfig() *Config {
	address := os.Getenv("SERVER_ADDRESS")
	if address == "" {
		address = consts.ServerDefaultAddress
	}
	return &Config{Server: ServerConfig{Address: address}}
}
