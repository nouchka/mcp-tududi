package config

import (
	"fmt"
	"os"
)

type Config struct {
	TududiAPIURL  string
	TududiAPIKey  string
	TududiEmail   string
	TududiPassword string
	HTTPPort      string
	LogLevel      string
}

func Load() Config {
	cfg := Config{
		TududiAPIURL:   getEnv("TUDUDI_API_URL", "http://localhost:3000"),
		TududiAPIKey:   getEnv("TUDUDI_API_KEY", ""),
		TududiEmail:    getEnv("TUDUDI_EMAIL", ""),
		TududiPassword: getEnv("TUDUDI_PASSWORD", ""),
		HTTPPort:       getEnv("HTTP_PORT", "8080"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
	}

	// Validate that at least one authentication method is provided
	if cfg.TududiAPIKey == "" && cfg.TududiEmail == "" {
		fmt.Fprintf(os.Stderr, "Warning: Neither TUDUDI_API_KEY nor TUDUDI_EMAIL is set\n")
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
