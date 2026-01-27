package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	ApiKey string
	Port   string
}

var (
	API_KEY string = "SECRET_API_KEY"

	// Default is 8888
	SERVER_PORT_KEY string = "PORT"
)

var (
	ErrApiKeyNotFound error = fmt.Errorf("API KEY not provided in .env")
)

func LoadConfig(filename string) (AppConfig, error) {
	env, err := godotenv.Read(filename)
	if err != nil {
		return AppConfig{}, err
	}

	apiKey, ok := env[API_KEY]
	if !ok {
		return AppConfig{}, ErrApiKeyNotFound
	}

	port, ok := env[SERVER_PORT_KEY]
	if !ok {
		port = ":8888"
	}

	return AppConfig{
		ApiKey: apiKey,
		Port:   port,
	}, nil
}
