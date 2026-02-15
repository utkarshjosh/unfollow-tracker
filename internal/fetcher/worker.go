package fetcher

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/utkarsh/unfollow-tracker/internal/config"
	"github.com/utkarsh/unfollow-tracker/internal/queue"
	"github.com/utkarsh/unfollow-tracker/internal/service"
)

// Worker processes fetch jobs from the queue
type Worker struct {
	id         int
	queue      *queue.Client
	instagram  *InstagramScraper
	fetcherSvc *service.FetcherService
	config     *config.Config
}

// NewWorker creates a new fetcher worker
func NewWorker(
	id int,
	queue *queue.Client,
	instagram *InstagramScraper,
	fetcherSvc *service.FetcherService,
	cfg *config.Config,
) *Worker {
	return &Worker{
		id:         id,
		queue:      queue,
		instagram:  instagram,
		fetcherSvc: fetcherSvc,
		config:     cfg,
	}
}

// Start begins processing jobs
func (w *Worker) Start(ctx context.Context) {
	log.Printf("Worker %d started", w.id)

	for {
		select {
		case <-ctx.Done():
			log.Printf("Worker %d stopping", w.id)
			return
		default:
			if err := w.processJob(ctx); err != nil {
				log.Printf("Worker %d error: %v", w.id, err)
				time.Sleep(5 * time.Second)
			}
		}
	}
}

func (w *Worker) processJob(ctx context.Context) error {
	// Dequeue job (blocking with 30 second timeout)
	job, err := w.queue.DequeueFetchJob(ctx, 30*time.Second)
	if err != nil {
		return fmt.Errorf("failed to dequeue: %w", err)
	}

	if job == nil {
		// Timeout, continue
		return nil
	}

	log.Printf("Worker %d processing account %s (%s)", w.id, job.Username, job.AccountID)

	// Parse account ID
	accountID, err := parseUUID(job.AccountID)
	if err != nil {
		return fmt.Errorf("invalid account ID: %w", err)
	}

	// TODO: Fetch from Instagram (stubbed for now)
	// For now, use placeholder data
	followerHashes := []string{
		"hash1", "hash2", "hash3", // Placeholder
	}

	// Process with service
	if err := w.fetcherSvc.ProcessFetchJob(ctx, accountID, job.ChunkID, followerHashes); err != nil {
		return fmt.Errorf("failed to process fetch job: %w", err)
	}

	log.Printf("Worker %d completed account %s", w.id, job.Username)
	return nil
}

func parseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}
