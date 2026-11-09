package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/singleflight"
)

var ErrCacheMiss = errors.New("key not found in cache or expired")

// CacheEntry wraps serialized data with an expiration timestamp.
type CacheEntry struct {
	Value     []byte
	ExpiresAt time.Time
}

// TTLCache provides thread-safe in-memory caching with expiration.
type TTLCache struct {
	mu      sync.RWMutex
	entries map[string]CacheEntry
}

func NewTTLCache() *TTLCache {
	return &TTLCache{
		entries: make(map[string]CacheEntry),
	}
}

func (c *TTLCache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.entries[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(entry.ExpiresAt) {
		return nil, false
	}

	return entry.Value, true
}

func (c *TTLCache) Set(key string, value []byte, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = CacheEntry{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}
}

func (c *TTLCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.entries, key)
}

// OrderRecord models cached order entity.
type OrderRecord struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Amount    float64   `json:"amount"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UnderlyingStore simulates a database repository.
type UnderlyingStore interface {
	FetchOrder(ctx context.Context, id string) (*OrderRecord, error)
	PersistOrder(ctx context.Context, o *OrderRecord) error
}

// CachedOrderService implements read-through caching with singleflight stampede defense.
type CachedOrderService struct {
	store    UnderlyingStore
	cache    *TTLCache
	ttl      time.Duration
	sfGroup  singleflight.Group
	dbReads  atomic.Int32
	dbWrites atomic.Int32
}

func NewCachedOrderService(store UnderlyingStore, cache *TTLCache, ttl time.Duration) *CachedOrderService {
	return &CachedOrderService{
		store: store,
		cache: cache,
		ttl:   ttl,
	}
}

// GetOrder implements read-through cache lookup.
func (s *CachedOrderService) GetOrder(ctx context.Context, id string) (*OrderRecord, error) {
	cacheKey := "order:" + id

	// 1. Check Cache
	if data, hit := s.cache.Get(cacheKey); hit {
		var o OrderRecord
		if err := json.Unmarshal(data, &o); err == nil {
			return &o, nil
		}
	}

	// 2. Cache Miss: Singleflight synchronization prevents cache stampedes
	v, err, _ := s.sfGroup.Do(cacheKey, func() (any, error) {
		// Double-check cache inside singleflight barrier in case a concurrent flight populated it
		if data, hit := s.cache.Get(cacheKey); hit {
			var o OrderRecord
			if err := json.Unmarshal(data, &o); err == nil {
				return &o, nil
			}
		}

		s.dbReads.Add(1)
		order, err := s.store.FetchOrder(ctx, id)
		if err != nil {
			return nil, err
		}
		if order == nil {
			return nil, nil
		}

		bytes, err := json.Marshal(order)
		if err == nil {
			s.cache.Set(cacheKey, bytes, s.ttl)
		}
		return order, nil
	})

	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}

	return v.(*OrderRecord), nil
}

// UpdateOrder modifies the resource and invalidates the cached key.
func (s *CachedOrderService) UpdateOrder(ctx context.Context, o *OrderRecord) error {
	s.dbWrites.Add(1)
	o.UpdatedAt = time.Now().UTC()

	// 1. Update source of truth in DB
	if err := s.store.PersistOrder(ctx, o); err != nil {
		return err
	}

	// 2. Invalidate cache key to prevent stale reads
	cacheKey := "order:" + o.ID
	s.cache.Delete(cacheKey)

	return nil
}

// MockDatabaseStore implements UnderlyingStore for tests.
type MockDatabaseStore struct {
	mu     sync.RWMutex
	orders map[string]*OrderRecord
}

func NewMockDatabaseStore() *MockDatabaseStore {
	return &MockDatabaseStore{
		orders: make(map[string]*OrderRecord),
	}
}

func (m *MockDatabaseStore) FetchOrder(ctx context.Context, id string) (*OrderRecord, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	o, exists := m.orders[id]
	if !exists {
		return nil, nil
	}
	// Return clone
	clone := *o
	return &clone, nil
}

func (m *MockDatabaseStore) PersistOrder(ctx context.Context, o *OrderRecord) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	clone := *o
	m.orders[o.ID] = &clone
	return nil
}

func main() {
	db := NewMockDatabaseStore()
	db.orders["ord-1"] = &OrderRecord{ID: "ord-1", Title: "Original Order", Amount: 50.0}

	cache := NewTTLCache()
	service := NewCachedOrderService(db, cache, 1*time.Minute)

	order, _ := service.GetOrder(context.Background(), "ord-1")
	fmt.Printf("Fetched order: %s ($%.2f)\n", order.Title, order.Amount)
}
