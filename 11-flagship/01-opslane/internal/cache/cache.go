// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

// Package cache provides an in-memory caching layer for the Opslane backend.
// It includes a Cache interface, TTL-based expiration, tenant-scoped key generation,
// write-through invalidation, and singleflight deduplication to prevent thundering herd.
package cache

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	// ErrNotFound (Error): returned when a cache key is not found
	ErrNotFound = errors.New("cache: key not found")
	// ErrKeyEmpty (Error): returned when an empty key is supplied
	ErrKeyEmpty = errors.New("cache: empty key")
	// ErrCacheClosed (Error): returned when the cache has been shut down
	ErrCacheClosed = errors.New("cache: closed")
)

// Cache (Interface): boundary that the rest of the application talks to for caching
type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	DeletePrefix(ctx context.Context, prefix string) error
	Close() error
}

// entry (Struct): holds one cached value and its absolute expiry time
type entry struct {
	value     []byte
	expiresAt time.Time
}

// expired (Method): returns true if the entry's expiry time is past
func (e entry) expired(now time.Time) bool {
	return !e.expiresAt.IsZero() && now.After(e.expiresAt)
}

// Config (Struct): controls the bounded behavior of the in-memory cache
type Config struct {
	MaxEntries int
	DefaultTTL time.Duration
}

// DefaultConfig (Function): returns production-reasonable cache defaults
func DefaultConfig() Config {
	return Config{
		MaxEntries: 4096,
		DefaultTTL: 5 * time.Minute,
	}
}

// TenantOrderKey (Function): builds a tenant-scoped cache key for a single order
func TenantOrderKey(tenantID, orderID int64) string {
	return fmt.Sprintf("tenant:%d:order:%d", tenantID, orderID)
}

// TenantOrderListKey (Function): builds a tenant-scoped cache key for the order list
func TenantOrderListKey(tenantID int64) string {
	return fmt.Sprintf("tenant:%d:orders", tenantID)
}

// TenantPaymentListKey (Function): builds a tenant-scoped cache key for payments by order
func TenantPaymentListKey(tenantID, orderID int64) string {
	return fmt.Sprintf("tenant:%d:order:%d:payments", tenantID, orderID)
}

// TenantOrderPrefix (Function): returns the prefix covering all order-related keys for a tenant
func TenantOrderPrefix(tenantID int64) string {
	return fmt.Sprintf("tenant:%d:order", tenantID)
}

// Invalidator (Struct): provides write-through cache invalidation for services
type Invalidator struct {
	cache Cache
}

// NewInvalidator (Constructor): creates an invalidator; nil cache makes all operations safe no-ops
func NewInvalidator(cache Cache) *Invalidator {
	return &Invalidator{cache: cache}
}

// InvalidateOrder (Method): removes cached data for a specific order and the tenant's order list
func (inv *Invalidator) InvalidateOrder(ctx context.Context, tenantID, orderID int64) {
	if inv == nil || inv.cache == nil {
		return
	}

	_ = inv.cache.Delete(ctx, TenantOrderKey(tenantID, orderID))
	_ = inv.cache.Delete(ctx, TenantOrderListKey(tenantID))
}

// InvalidatePayments (Method): removes cached payment data for an order
func (inv *Invalidator) InvalidatePayments(ctx context.Context, tenantID, orderID int64) {
	if inv == nil || inv.cache == nil {
		return
	}

	_ = inv.cache.Delete(ctx, TenantPaymentListKey(tenantID, orderID))
}

// InvalidateTenantOrders (Method): removes all cached order and payment data for a tenant
func (inv *Invalidator) InvalidateTenantOrders(ctx context.Context, tenantID int64) {
	if inv == nil || inv.cache == nil {
		return
	}

	_ = inv.cache.DeletePrefix(ctx, TenantOrderPrefix(tenantID))
	_ = inv.cache.Delete(ctx, TenantOrderListKey(tenantID))
}

// NoopCache (Struct): satisfies the Cache interface without storing anything; useful for testing
type NoopCache struct{}

// Get (Method): returns ErrNotFound for the no-op cache
func (NoopCache) Get(context.Context, string) ([]byte, error) { return nil, ErrNotFound }

// Set (Method): no-op for the no-op cache
func (NoopCache) Set(context.Context, string, []byte, time.Duration) error { return nil }

// Delete (Method): no-op for the no-op cache
func (NoopCache) Delete(context.Context, string) error { return nil }

// DeletePrefix (Method): no-op for the no-op cache
func (NoopCache) DeletePrefix(context.Context, string) error { return nil }

// Close (Method): no-op for the no-op cache
func (NoopCache) Close() error { return nil }

// compile-time interface checks
var (
	_ Cache = (*NoopCache)(nil)
	_ Cache = (*InMemoryStore)(nil)
)

// Singleflight (Struct): deduplicates concurrent cache loads for the same key to prevent thundering herd
type Singleflight struct {
	mu    sync.Mutex
	calls map[string]*call
}

// call (Struct): tracks a single in-flight singleflight operation
type call struct {
	wg  sync.WaitGroup
	val []byte
	err error
}

// Do (Method): executes fn once per key, returning the same result to concurrent callers
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
	// instead of nil, nil - which would look like a successful empty result.
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
