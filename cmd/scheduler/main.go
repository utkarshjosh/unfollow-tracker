package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/utkarsh/unfollow-tracker/internal/config"
	"github.com/utkarsh/unfollow-tracker/internal/database"
	"github.com/utkarsh/unfollow-tracker/internal/queue"
	"github.com/utkarsh/unfollow-tracker/internal/repository"
	"github.com/utkarsh/unfollow-tracker/internal/scheduler"
)

func main() {
	// Load .env file if present
	_ = godotenv.Load()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := database.NewDatabase(cfg.GetDatabaseConfig())
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close(db)

	// Initialize queue
	q, err := queue.NewClient(cfg.Redis.URL)
	if err != nil {
		log.Fatalf("Failed to connect to queue: %v", err)
	}

	// Initialize repository
	accountRepo := repository.NewPostgresAccountRepository(db)

	// Create scheduler
	sched := scheduler.NewScheduler(accountRepo, q, cfg)

	// Start scheduler
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go sched.Start(ctx)

	// Wait for interrupt
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down scheduler...")
	cancel()
}
