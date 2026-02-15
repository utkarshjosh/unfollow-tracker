package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	// Queue names
	FetchJobQueue     = "queue:fetch_jobs"
	NotificationQueue = "queue:notifications"

	// Key prefixes
	RateLimitPrefix = "ratelimit:"
	CachePrefix     = "cache:"
)

// Client wraps Redis operations for the queue
type Client struct {
	rdb *redis.Client
}

// NewClient creates a new queue client
func NewClient(redisURL string) (*Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("invalid redis URL: %w", err)
	}

	rdb := redis.NewClient(opts)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	return &Client{rdb: rdb}, nil
}

// Close closes the Redis connection
func (c *Client) Close() error {
	return c.rdb.Close()
}

// FetchJob represents a job to fetch follower data
type FetchJob struct {
	ID        string    `json:"id"`
	AccountID string    `json:"account_id"`
	ChunkID   int       `json:"chunk_id"`
	Username  string    `json:"username"`
	Platform  string    `json:"platform"`
	Priority  int       `json:"priority"`
	CreatedAt time.Time `json:"created_at"`
}

// EnqueueFetchJob adds a fetch job to the queue
func (c *Client) EnqueueFetchJob(ctx context.Context, job *FetchJob) error {
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}

	return c.rdb.LPush(ctx, FetchJobQueue, data).Err()
}

// DequeueFetchJob retrieves the next fetch job (blocking)
func (c *Client) DequeueFetchJob(ctx context.Context, timeout time.Duration) (*FetchJob, error) {
	result, err := c.rdb.BRPop(ctx, timeout, FetchJobQueue).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // No job available
		}
		return nil, err
	}

	if len(result) < 2 {
		return nil, nil
	}

	var job FetchJob
	if err := json.Unmarshal([]byte(result[1]), &job); err != nil {
		return nil, err
	}

	return &job, nil
}

// QueueLength returns the number of pending jobs
func (c *Client) QueueLength(ctx context.Context, queue string) (int64, error) {
	return c.rdb.LLen(ctx, queue).Result()
}

// SetFollowerHashes stores current follower hashes for quick comparison
func (c *Client) SetFollowerHashes(ctx context.Context, accountID string, chunkID int, hashes []string) error {
	key := fmt.Sprintf("followers:%s:%d", accountID, chunkID)

	// Delete existing and add new
	pipe := c.rdb.Pipeline()
	pipe.Del(ctx, key)
	if len(hashes) > 0 {
		members := make([]interface{}, len(hashes))
		for i, h := range hashes {
			members[i] = h
		}
		pipe.SAdd(ctx, key, members...)
		pipe.Expire(ctx, key, 48*time.Hour) // TTL for cache
	}

	_, err := pipe.Exec(ctx)
	return err
}

// GetFollowerHashes retrieves cached follower hashes
func (c *Client) GetFollowerHashes(ctx context.Context, accountID string, chunkID int) ([]string, error) {
	key := fmt.Sprintf("followers:%s:%d", accountID, chunkID)
	return c.rdb.SMembers(ctx, key).Result()
}

// CheckRateLimit checks if an action is rate limited
func (c *Client) CheckRateLimit(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	fullKey := RateLimitPrefix + key

	count, err := c.rdb.Incr(ctx, fullKey).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		c.rdb.Expire(ctx, fullKey, window)
	}

	return count <= int64(limit), nil
}
