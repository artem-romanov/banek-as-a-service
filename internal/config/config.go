package config

import (
	"fmt"
	"log/slog"
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
	if s == "" {
		slog.Warn("Environment variable not set, using default (dev)")
		return EnvDev
	}

	lowerS := strings.ToLower(s)
	env, ok := strToEnv[lowerS]
	if !ok {
		slog.Warn("Environment variable not found, using default (dev)", "key", s)
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

	CerebrusToken string

	GroqToken string
}

var (
	API_KEY string = "SECRET_API_KEY"

	// Possible options: dev, prod
	ENV_KEY string = "ENV"

	SERVER_PORT_KEY string = "PORT"

	CEREBRUS_TOKEN string = "CEREBRUS_TOKEN"

	GROQ_TOKEN string = "GROQ_TOKEN"
)

func LoadConfig(filename string) (AppConfig, error) {
	env, err := godotenv.Read(filename)
	if err != nil {
		return AppConfig{}, err
	}

	apiKey, ok := env[API_KEY]
	if !ok {
		return AppConfig{}, fmt.Errorf("API KEY not provided in .env")
	}

	// by default env[ENV_KEY] == ""
	// additional check if !ok is not needed
	e := ParseEnvType(env[ENV_KEY])

	port, ok := env[SERVER_PORT_KEY]
	if !ok {
		slog.Warn("Port is not set. Using default 8888.")
		port = "8888"
	}

	cerebrusToken, ok := env[CEREBRUS_TOKEN]
	if !ok {
		slog.Warn("Cerebrus token is not set in env. AI might wont work!")
		cerebrusToken = ""
	}

	groqToken, ok := env[GROQ_TOKEN]
	if !ok {
		slog.Warn("Groq token is not set in env. AI might wont work!")
		cerebrusToken = ""
	}

	return AppConfig{
		ApiKey:        apiKey,
		Port:          port,
		Environment:   e,
		CerebrusToken: cerebrusToken,
		GroqToken:     groqToken,
	}, nil
}
