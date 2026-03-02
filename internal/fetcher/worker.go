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

	// Check if we have session authentication
	if !w.instagram.IsAuthenticated() {
		return fmt.Errorf("no Instagram session cookie configured - cannot fetch followers")
	}

	// Fetch profile to get user ID and follower count
	profile, err := w.instagram.FetchProfile(ctx, job.Username)
	if err != nil {
		return fmt.Errorf("failed to fetch profile for %s: %w", job.Username, err)
	}

	log.Printf("Worker %d: %s has %d followers (public: %v)", w.id, profile.Username, profile.FollowerCount, profile.IsPublic)

	// If account is private and we're not following, we can't fetch followers
	if !profile.IsPublic {
		return fmt.Errorf("account %s is private - cannot fetch followers without following", job.Username)
	}

	// Fetch all followers with pagination
	delay := time.Duration(w.config.Scraper.DelayMs) * time.Millisecond
	allFollowers, err := w.instagram.FetchAllFollowers(ctx, profile.UserID, delay)
	if err != nil {
		return fmt.Errorf("failed to fetch followers for %s: %w", job.Username, err)
	}

	log.Printf("Worker %d: fetched %d followers for %s", w.id, len(allFollowers), job.Username)

	// Chunk followers (1000 per chunk)
	chunks := ChunkFollowers(allFollowers, 1000)

	// Process each chunk
	for chunkIdx, chunk := range chunks {
		if err := w.fetcherSvc.ProcessFetchJob(ctx, accountID, chunkIdx, chunk); err != nil {
			return fmt.Errorf("failed to process chunk %d: %w", chunkIdx, err)
		}

		// Apply delay between chunks to avoid rate limiting
		if delay > 0 && chunkIdx < len(chunks)-1 {
			time.Sleep(delay)
		}
	}

	log.Printf("Worker %d completed account %s (%d followers in %d chunks)", w.id, job.Username, len(allFollowers), len(chunks))
	return nil
}

func parseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}
