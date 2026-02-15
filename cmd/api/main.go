package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/utkarsh/unfollow-tracker/internal/api"
	"github.com/utkarsh/unfollow-tracker/internal/config"
	"github.com/utkarsh/unfollow-tracker/internal/database"
)

func main() {
	// Load .env file if present
	_ = godotenv.Load()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database connection
	db, err := database.NewDatabase(cfg.GetDatabaseConfig())
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close(db)

	// Create and start server
	server := api.NewServer(cfg, db)
	if err := server.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
