package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
		check   func(*testing.T, Config)
	}{
		{
			name:    "default values",
			envVars: map[string]string{},
			check: func(t *testing.T, cfg Config) {
				if cfg.TududiAPIURL != "http://localhost:3000" {
					t.Errorf("expected default TududiAPIURL, got %s", cfg.TududiAPIURL)
				}
				if cfg.HTTPPort != "8080" {
					t.Errorf("expected default HTTPPort, got %s", cfg.HTTPPort)
				}
				if cfg.LogLevel != "info" {
					t.Errorf("expected default LogLevel, got %s", cfg.LogLevel)
				}
			},
		},
		{
			name: "custom values",
			envVars: map[string]string{
				"TUDUDI_API_URL": "http://example.com:3000",
				"TUDUDI_API_KEY": "test-key",
				"HTTP_PORT":      "9000",
				"LOG_LEVEL":      "debug",
			},
			check: func(t *testing.T, cfg Config) {
				if cfg.TududiAPIURL != "http://example.com:3000" {
					t.Errorf("expected custom TududiAPIURL, got %s", cfg.TududiAPIURL)
				}
				if cfg.TududiAPIKey != "test-key" {
					t.Errorf("expected custom TududiAPIKey, got %s", cfg.TududiAPIKey)
				}
				if cfg.HTTPPort != "9000" {
					t.Errorf("expected custom HTTPPort, got %s", cfg.HTTPPort)
				}
				if cfg.LogLevel != "debug" {
					t.Errorf("expected custom LogLevel, got %s", cfg.LogLevel)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear env vars
			os.Clearenv()

			// Set test env vars
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			cfg := Load()
			tt.check(t, cfg)
		})
	}
}
