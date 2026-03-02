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
	"github.com/utkarsh/unfollow-tracker/internal/fetcher"
	"github.com/utkarsh/unfollow-tracker/internal/queue"
	"github.com/utkarsh/unfollow-tracker/internal/repository"
	"github.com/utkarsh/unfollow-tracker/internal/service"
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

	// Initialize repositories
	accountRepo := repository.NewPostgresAccountRepository(db)
	snapshotRepo := repository.NewPostgresSnapshotRepository(db)
	unfollowRepo := repository.NewPostgresUnfollowRepository(db)

	// Initialize services
	fetcherSvc := service.NewFetcherService(accountRepo, snapshotRepo, unfollowRepo)

	// Create Instagram scraper with session cookie
	// TODO: Load proxies from config if available
	proxyPool := fetcher.NewProxyPool([]string{})
	instagram := fetcher.NewInstagramScraper(proxyPool, cfg.Scraper.InstagramSession)

	// Create workers
	workers := make([]*fetcher.Worker, cfg.Scraper.MaxConcurrent)
	for i := 0; i < cfg.Scraper.MaxConcurrent; i++ {
		workers[i] = fetcher.NewWorker(i+1, q, instagram, fetcherSvc, cfg)
	}

	// Start workers
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, w := range workers {
		go w.Start(ctx)
	}

	// Wait for interrupt
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down workers...")
	cancel()
}
