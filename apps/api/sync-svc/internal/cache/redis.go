// Package cache wraps go-redis for sync-svc caching (Valkey-compatible).
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const DefaultTTL = 30 * time.Minute

type Client struct {
	rdb *redis.Client
}

func New(redisURL string) (*Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}
	return &Client{rdb: redis.NewClient(opts)}, nil
}

func (c *Client) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}

// Set marshals payload to JSON and stores it with ttl.
func (c *Client) Set(ctx context.Context, key string, payload any, ttl time.Duration) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	return c.rdb.Set(ctx, key, data, ttl).Err()
}

// Get retrieves and unmarshals a cached value. Returns (false, nil) on cache miss.
func (c *Client) Get(ctx context.Context, key string, dest any) (bool, error) {
	data, err := c.rdb.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("redis get: %w", err)
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return false, fmt.Errorf("unmarshal: %w", err)
	}
	return true, nil
}

// Del deletes one or more keys.
func (c *Client) Del(ctx context.Context, keys ...string) error {
	return c.rdb.Del(ctx, keys...).Err()
}

// Keys returns matching key names (for cache invalidation).
func (c *Client) Keys(ctx context.Context, pattern string) ([]string, error) {
	return c.rdb.Keys(ctx, pattern).Result()
}

// Key helpers for consistent naming
func LoanKey(borrowerID string) string { return "sync:loans:" + borrowerID }
func ScheduleKey(loanID string) string { return "sync:schedule:" + loanID }
func DebtsKey(loanID string) string    { return "sync:debts:" + loanID }
func AllLoansKey() string              { return "sync:all_loans" }
