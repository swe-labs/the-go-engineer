// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Package ratelimit provides distributed rate limiting using PostgreSQL.
// Role: Rate limiting boundary - prevents service overload from coordinated clients.
// Boundary: Rate limit state is shared across all application instances.
// Failure mode: If DB is unavailable, fail open (allow request) to prevent outages.

package ratelimit

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

type Config struct {
	RequestsPerSecond int
	BurstSize         int
	WindowSeconds     int
	DB                *sql.DB
}

type Limiter struct {
	cfg Config

	mu       sync.Mutex
	counters map[string]*windowCounter
}

type Decision struct {
	Allowed   bool
	Limit     int
	Remaining int
	ResetAt   time.Time
}

type windowCounter struct {
	count     int64
	windowEnd time.Time
}

func New(cfg Config) *Limiter {
	if cfg.RequestsPerSecond == 0 {
		cfg.RequestsPerSecond = 10
	}
	if cfg.BurstSize == 0 {
		cfg.BurstSize = 20
	}
	if cfg.WindowSeconds == 0 {
		cfg.WindowSeconds = 1
	}

	return &Limiter{
		cfg:      cfg,
		counters: make(map[string]*windowCounter),
	}
}

func (l *Limiter) Allow(key string) (bool, error) {
	d, err := l.AllowWithDecision(context.Background(), key)
	return d.Allowed, err
}

func (l *Limiter) AllowContext(ctx context.Context, key string) (bool, error) {
	d, err := l.AllowWithDecision(ctx, key)
	return d.Allowed, err
}

func (l *Limiter) AllowWithDecision(ctx context.Context, key string) (Decision, error) {
	if l.cfg.DB == nil {
		return l.localAllow(key), nil
	}
	return l.dbAllow(ctx, key)
}

func (l *Limiter) localAllow(key string) Decision {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	counter, exists := l.counters[key]

	if !exists || now.After(counter.windowEnd) {
		l.counters[key] = &windowCounter{
			count:     1,
			windowEnd: now.Add(time.Duration(l.cfg.WindowSeconds) * time.Second),
		}
		return Decision{
			Allowed:   true,
			Limit:     l.cfg.BurstSize,
			Remaining: l.cfg.BurstSize - 1,
			ResetAt:   l.counters[key].windowEnd,
		}
	}

	if counter.count >= int64(l.cfg.BurstSize) {
		return Decision{
			Allowed:   false,
			Limit:     l.cfg.BurstSize,
			Remaining: 0,
			ResetAt:   counter.windowEnd,
		}
	}

	counter.count++
	remaining := l.cfg.BurstSize - int(counter.count)
	if remaining < 0 {
		remaining = 0
	}
	return Decision{
		Allowed:   true,
		Limit:     l.cfg.BurstSize,
		Remaining: remaining,
		ResetAt:   counter.windowEnd,
	}
}

func (l *Limiter) dbAllow(ctx context.Context, key string) (Decision, error) {
	now := time.Now()
	window := now.Truncate(time.Duration(l.cfg.WindowSeconds) * time.Second)
	windowEnd := window.Add(time.Duration(l.cfg.WindowSeconds) * time.Second)

	var count int64
	err := l.cfg.DB.QueryRowContext(ctx, `
		INSERT INTO rate_limits (key, window, count, expires_at)
		VALUES ($1, $2, 1, $3)
		ON CONFLICT (key, window) DO UPDATE SET
			count = rate_limits.count + 1,
			expires_at = $3
		RETURNING (
			SELECT count FROM rate_limits WHERE key = $1 AND window = $2
		)
	`, key, window, windowEnd).Scan(&count)

	if err != nil {
		if err == sql.ErrNoRows {
			return Decision{
				Allowed:   true,
				Limit:     l.cfg.BurstSize,
				Remaining: l.cfg.BurstSize - 1,
				ResetAt:   windowEnd,
			}, nil
		}
		return Decision{
			Allowed:   true,
			Limit:     l.cfg.BurstSize,
			Remaining: l.cfg.BurstSize,
			ResetAt:   windowEnd,
		}, err
	}

	remaining := l.cfg.BurstSize - int(count)
	if remaining < 0 {
		remaining = 0
	}
	return Decision{
		Allowed:   count <= int64(l.cfg.BurstSize),
		Limit:     l.cfg.BurstSize,
		Remaining: remaining,
		ResetAt:   windowEnd,
	}, nil
}

func (l *Limiter) Reset(key string) error {
	if l.cfg.DB == nil {
		l.mu.Lock()
		defer l.mu.Unlock()
		delete(l.counters, key)
		return nil
	}

	_, err := l.cfg.DB.Exec("DELETE FROM rate_limits WHERE key = $1", key)
	return err
}

type MultiLimiter struct {
	limiters map[string]*Limiter
	mu       sync.Mutex
}

func NewMultiLimiter() *MultiLimiter {
	return &MultiLimiter{
		limiters: make(map[string]*Limiter),
	}
}

func (m *MultiLimiter) Register(name string, cfg Config) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.limiters[name] = New(cfg)
}

func (m *MultiLimiter) Allow(name, key string) (bool, error) {
	m.mu.Lock()
	limiter, ok := m.limiters[name]
	m.mu.Unlock()

	if !ok {
		return true, nil
	}

	return limiter.Allow(key)
}

func CreateRateLimitTable(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS rate_limits (
			key TEXT NOT NULL,
			window TIMESTAMPTZ NOT NULL,
			count BIGINT NOT NULL DEFAULT 1,
			expires_at TIMESTAMPTZ NOT NULL,
			PRIMARY KEY (key, window)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create rate_limits table: %w", err)
	}

	_, err = db.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS idx_rate_limits_expires
		ON rate_limits (expires_at)
	`)
	return err
}
