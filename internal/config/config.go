package config

import (
	"fmt"
	"strings"

	"github.com/joho/godotenv"
)

type EnvType int

const (
	EnvDev EnvType = iota
	EnvProd
)

var strToEnv = map[string]EnvType{
	"dev":  EnvDev,
	"prod": EnvProd,
}

// ParseEnvType returns EnvType based on strToEnv map.
// String lowercased for consistency.
// It defaults to dev if string is not found.
func ParseEnvType(s string) EnvType {
	lowerS := strings.ToLower(s)

	env, ok := strToEnv[lowerS]
	if !ok {
		return EnvDev
	}
	return env
}

// AppConfig reporesents application configuration
type AppConfig struct {
	// Key which protects all endpoints
	// For security reasons and simplicity
	ApiKey string

	// Default is 8888
	Port string

	// Application environment. Defaults to dev.
	Environment EnvType
}

var (
	API_KEY string = "SECRET_API_KEY"

	// Possible options: dev, prod
	ENV_KEY string = "ENV"

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

	// by default env[ENV_KEY] == ""
	// additional check if !ok is not needed
	e := ParseEnvType(env[ENV_KEY])

	port, ok := env[SERVER_PORT_KEY]
	if !ok {
		port = "8888"
	}

	return AppConfig{
		ApiKey:      apiKey,
		Port:        port,
		Environment: e,
	}, nil
}
