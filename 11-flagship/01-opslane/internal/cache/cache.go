// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package cache

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrNotFound    = errors.New("cache: key not found")
	ErrKeyEmpty    = errors.New("cache: empty key")
	ErrCacheClosed = errors.New("cache: closed")
)

// Cache is the boundary that the rest of the application talks to.
// PostgreSQL stays the system of record; this is always additive.
type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	DeletePrefix(ctx context.Context, prefix string) error
	Close() error
}

// entry holds one cached value and its absolute expiry time.
type entry struct {
	value     []byte
	expiresAt time.Time
}

func (e entry) expired(now time.Time) bool {
	return !e.expiresAt.IsZero() && now.After(e.expiresAt)
}

// Config controls the bounded behavior of the in-memory cache.
type Config struct {
	MaxEntries int
	DefaultTTL time.Duration
}

// DefaultConfig returns production-reasonable cache defaults.
func DefaultConfig() Config {
	return Config{
		MaxEntries: 4096,
		DefaultTTL: 5 * time.Minute,
	}
}

// TenantOrderKey builds a tenant-scoped cache key for a single order.
func TenantOrderKey(tenantID, orderID int64) string {
	return fmt.Sprintf("tenant:%d:order:%d", tenantID, orderID)
}

// TenantOrderListKey builds a tenant-scoped cache key for the order list.
func TenantOrderListKey(tenantID int64) string {
	return fmt.Sprintf("tenant:%d:orders", tenantID)
}

// TenantPaymentListKey builds a tenant-scoped cache key for payments by order.
func TenantPaymentListKey(tenantID, orderID int64) string {
	return fmt.Sprintf("tenant:%d:order:%d:payments", tenantID, orderID)
}

// TenantOrderPrefix returns the prefix that covers all order-related keys
// for a given tenant, so that a write can invalidate the whole group.
func TenantOrderPrefix(tenantID int64) string {
	return fmt.Sprintf("tenant:%d:order", tenantID)
}

// Invalidator provides write-through cache invalidation.
// Services call these methods after mutations so the cache never
// silently serves stale data.
type Invalidator struct {
	cache Cache
}

// NewInvalidator creates an invalidator. If cache is nil, all operations
// are safe no-ops.
func NewInvalidator(cache Cache) *Invalidator {
	return &Invalidator{cache: cache}
}

// InvalidateOrder removes cached data for a specific order and the
// tenant's order list so the next read comes from PostgreSQL.
func (inv *Invalidator) InvalidateOrder(ctx context.Context, tenantID, orderID int64) {
	if inv == nil || inv.cache == nil {
		return
	}

	_ = inv.cache.Delete(ctx, TenantOrderKey(tenantID, orderID))
	_ = inv.cache.Delete(ctx, TenantOrderListKey(tenantID))
}

// InvalidatePayments removes cached payment data for an order.
func (inv *Invalidator) InvalidatePayments(ctx context.Context, tenantID, orderID int64) {
	if inv == nil || inv.cache == nil {
		return
	}

	_ = inv.cache.Delete(ctx, TenantPaymentListKey(tenantID, orderID))
}

// InvalidateTenantOrders removes all cached order and payment data for a
// tenant. Use this after bulk mutations or when granular invalidation
// would be too expensive.
func (inv *Invalidator) InvalidateTenantOrders(ctx context.Context, tenantID int64) {
	if inv == nil || inv.cache == nil {
		return
	}

	_ = inv.cache.DeletePrefix(ctx, TenantOrderPrefix(tenantID))
	_ = inv.cache.Delete(ctx, TenantOrderListKey(tenantID))
}

// NoopCache satisfies the Cache interface without storing anything.
// Useful for testing and for environments where caching is disabled.
type NoopCache struct{}

func (NoopCache) Get(context.Context, string) ([]byte, error)              { return nil, ErrNotFound }
func (NoopCache) Set(context.Context, string, []byte, time.Duration) error { return nil }
func (NoopCache) Delete(context.Context, string) error                     { return nil }
func (NoopCache) DeletePrefix(context.Context, string) error               { return nil }
func (NoopCache) Close() error                                             { return nil }

// compile-time interface checks
var (
	_ Cache = (*NoopCache)(nil)
	_ Cache = (*InMemoryStore)(nil)
)

// Singleflight deduplicates concurrent loads for the same key.
// Callers wrap their database read inside fn; if multiple goroutines
// request the same key before the first load completes, only one
// actually calls fn and the rest receive its result.
type Singleflight struct {
	mu    sync.Mutex
	calls map[string]*call
}

type call struct {
	wg  sync.WaitGroup
	val []byte
	err error
}

// Do executes fn once per key. Concurrent callers with the same key
// block until the first caller's fn returns, then receive the same
// result.
func (sf *Singleflight) Do(key string, fn func() ([]byte, error)) (val []byte, err error) {
	sf.mu.Lock()
	if sf.calls == nil {
		sf.calls = make(map[string]*call)
	}

	if c, ok := sf.calls[key]; ok {
		sf.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}

	c := &call{}
	c.wg.Add(1)
	sf.calls[key] = c
	sf.mu.Unlock()

	// Ensure cleanup runs even if fn() panics.
	// If a panic occurs, set c.err so that waiters receive an error
	// instead of nil, nil — which would look like a successful empty result.
	// Named returns ensure the leader also receives the error.
	defer func() {
		if r := recover(); r != nil {
			c.err = fmt.Errorf("cache: singleflight panic: %v", r)
			err = c.err
		}
		c.wg.Done()
		sf.mu.Lock()
		delete(sf.calls, key)
		sf.mu.Unlock()
	}()

	c.val, c.err = fn()
	val = c.val
	err = c.err
	return
}
