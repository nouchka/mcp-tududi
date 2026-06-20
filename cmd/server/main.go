package main

import (
	"log"
	"os"

	"github.com/nouchka/mcp-tududi/internal/client"
	"github.com/nouchka/mcp-tududi/internal/config"
	"github.com/nouchka/mcp-tududi/internal/mcp"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Create Tududi client
	tududuClient := client.NewTududuClient(
		cfg.TududiAPIURL,
		cfg.TududiAPIKey,
		cfg.TududiEmail,
		cfg.TududiPassword,
	)

	// Authenticate if needed
	if err := tududuClient.Authenticate(); err != nil {
		log.Printf("Warning: authentication error: %v", err)
	}

	// Create and start MCP server
	server := mcp.NewServer(cfg.HTTPPort, tududuClient)
	log.Printf("Starting MCP server on port %s", cfg.HTTPPort)
	log.Printf("Tududi API URL: %s", cfg.TududiAPIURL)

	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
		os.Exit(1)
	}
}
