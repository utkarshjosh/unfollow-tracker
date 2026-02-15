package scheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/utkarsh/unfollow-tracker/internal/config"
	"github.com/utkarsh/unfollow-tracker/internal/queue"
	"github.com/utkarsh/unfollow-tracker/internal/repository"
)

// Scheduler manages periodic scan jobs
type Scheduler struct {
	accountRepo repository.AccountRepository
	queue       *queue.Client
	config      *config.Config
}

// NewScheduler creates a new scheduler
func NewScheduler(
	accountRepo repository.AccountRepository,
	queue *queue.Client,
	cfg *config.Config,
) *Scheduler {
	return &Scheduler{
		accountRepo: accountRepo,
		queue:       queue,
		config:      cfg,
	}
}

// Start begins the scheduler
func (s *Scheduler) Start(ctx context.Context) error {
	log.Println("Scheduler starting...")

	// Run different scheduling loops
	go s.accountScanLoop(ctx)
	go s.cleanupLoop(ctx)
	go s.notificationLoop(ctx)

	<-ctx.Done()
	return nil
}

// accountScanLoop checks for accounts that need scanning
func (s *Scheduler) accountScanLoop(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := s.enqueueScans(ctx); err != nil {
				log.Printf("Error enqueuing scans: %v", err)
			}
		}
	}
}

func (s *Scheduler) enqueueScans(ctx context.Context) error {
	// Find accounts needing scan (interval of 24 hours, limit 100)
	accounts, err := s.accountRepo.FindAccountsNeedingScan(ctx, 24, 100)
	if err != nil {
		return fmt.Errorf("failed to find accounts: %w", err)
	}

	if len(accounts) == 0 {
		return nil
	}

	log.Printf("Enqueueing %d accounts for scanning", len(accounts))

	// Enqueue each account
	for _, account := range accounts {
		job := &queue.FetchJob{
			AccountID: account.ID.String(),
			Platform:  string(account.Platform),
			Username:  account.Username,
			ChunkID:   0, // Start with chunk 0
		}

		if err := s.queue.EnqueueFetchJob(ctx, job); err != nil {
			log.Printf("Failed to enqueue account %s: %v", account.ID, err)
			continue
		}
	}

	return nil
}

// cleanupLoop performs periodic cleanup tasks
func (s *Scheduler) cleanupLoop(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.performCleanup(ctx)
		}
	}
}

func (s *Scheduler) performCleanup(ctx context.Context) {
	// TODO: Implement
	// 1. Delete old snapshots (keep last N per account)
	// 2. Aggregate old unfollows into summaries
	// 3. Clean up stale queue jobs

	log.Println("Performing cleanup...")
}

// notificationLoop handles delayed notifications
func (s *Scheduler) notificationLoop(ctx context.Context) {
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.processNotifications(ctx)
		}
	}
}

func (s *Scheduler) processNotifications(ctx context.Context) {
	// TODO: Implement
	// 1. Find unfollows that are old enough to notify (delay for ethical reasons)
	// 2. Aggregate into summaries (not individual alerts)
	// 3. Send notifications (email, push)
	// 4. Mark as notified

	log.Println("Processing notifications...")
}
