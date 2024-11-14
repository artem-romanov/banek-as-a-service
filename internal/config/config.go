package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	ApiKey string
}

var (
	API_KEY string = "SECRET_API_KEY"
)

var (
	ERR_API_KEY_NOT_FOUND error = fmt.Errorf("API KEY not provided in .env")
)

func LoadConfig(filename string) (AppConfig, error) {
	env, err := godotenv.Read(filename)
	if err != nil {
		return AppConfig{}, err
	}

	apiKey, ok := env[API_KEY]
	if !ok {
		return AppConfig{}, ERR_API_KEY_NOT_FOUND
	}

	return AppConfig{
		ApiKey: apiKey,
	}, nil
}
