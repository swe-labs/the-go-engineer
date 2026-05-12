// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package cache

import (
	"context"
	"strings"
	"sync"
	"time"
)

// InMemoryStore (Struct): bounded, TTL-aware in-memory cache with concurrent access support
//
// Design decisions:
//   - sync.RWMutex protects the map; reads take the read-lock, writes take
//     the write-lock. This keeps contention low under read-heavy workloads.
//   - When the cache reaches MaxEntries, the oldest entry is evicted. This
//     is not true LRU (no access-time tracking) - it is insert-order eviction.
//     For a teaching codebase this is an acceptable simplification; real
//     production caches use probabilistic LRU or slab allocators.
//   - Expired entries are cleaned lazily on Get and periodically by a
//     background goroutine. The background janitor prevents slow memory
//     growth when keys are set but never read again.
//   - The store is safe to use concurrently from multiple goroutines.
type InMemoryStore struct {
	mu      sync.RWMutex
	entries map[string]entry
	order   []string // insertion order for eviction
	config  Config
	now     func() time.Time

	closed chan struct{}
	once   sync.Once
}

// NewInMemoryStore (Constructor): creates a bounded in-memory cache with a background cleanup janitor
func NewInMemoryStore(config Config) *InMemoryStore {
	if config.MaxEntries <= 0 {
		config.MaxEntries = DefaultConfig().MaxEntries
	}
	if config.DefaultTTL <= 0 {
		config.DefaultTTL = DefaultConfig().DefaultTTL
	}

	s := &InMemoryStore{
		entries: make(map[string]entry, config.MaxEntries),
		order:   make([]string, 0, config.MaxEntries),
		config:  config,
		now:     time.Now,
		closed:  make(chan struct{}),
	}

	go s.janitor(30 * time.Second)

	return s
}

// Get (Method): retrieves a value by key, returning a copy to prevent mutation
func (s *InMemoryStore) Get(_ context.Context, key string) ([]byte, error) {
	if key == "" {
		return nil, ErrKeyEmpty
	}

	select {
	case <-s.closed:
		return nil, ErrCacheClosed
	default:
	}

	s.mu.RLock()
	e, ok := s.entries[key]
	s.mu.RUnlock()

	if !ok {
		return nil, ErrNotFound
	}

	if e.expired(s.now()) {
		// Lazy eviction: delete expired entry on read miss.
		s.mu.Lock()
		if current, stillThere := s.entries[key]; stillThere && current.expired(s.now()) {
			delete(s.entries, key)
			s.removeFromOrder(key)
		}
		s.mu.Unlock()
		return nil, ErrNotFound
	}

	// Return a copy so the caller cannot mutate cached data.
	cp := make([]byte, len(e.value))
	copy(cp, e.value)
	return cp, nil
}

// Set (Method): stores a value with TTL, evicting oldest entry if at capacity
func (s *InMemoryStore) Set(_ context.Context, key string, value []byte, ttl time.Duration) error {
	if key == "" {
		return ErrKeyEmpty
	}

	select {
	case <-s.closed:
		return ErrCacheClosed
	default:
	}

	if ttl <= 0 {
		ttl = s.config.DefaultTTL
	}

	// Copy value so the caller cannot mutate cached data after Set returns.
	cp := make([]byte, len(value))
	copy(cp, value)

	s.mu.Lock()
	defer s.mu.Unlock()

	// If the key already exists, update in place without changing order.
	if _, exists := s.entries[key]; exists {
		s.entries[key] = entry{value: cp, expiresAt: s.now().Add(ttl)}
		return nil
	}

	// Evict the oldest entry if at capacity.
	if len(s.entries) >= s.config.MaxEntries {
		s.evictOldest()
	}

	s.entries[key] = entry{value: cp, expiresAt: s.now().Add(ttl)}
	s.order = append(s.order, key)
	return nil
}

// Delete (Method): removes a single key from the cache
func (s *InMemoryStore) Delete(_ context.Context, key string) error {
	if key == "" {
		return ErrKeyEmpty
	}

	select {
	case <-s.closed:
		return ErrCacheClosed
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.entries[key]; exists {
		delete(s.entries, key)
		s.removeFromOrder(key)
	}
	return nil
}

// DeletePrefix (Method): removes all entries whose key starts with the given prefix
func (s *InMemoryStore) DeletePrefix(_ context.Context, prefix string) error {
	if prefix == "" {
		return ErrKeyEmpty
	}

	select {
	case <-s.closed:
		return ErrCacheClosed
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	var toRemove []string
	for key := range s.entries {
		if strings.HasPrefix(key, prefix) {
			toRemove = append(toRemove, key)
		}
	}

	for _, key := range toRemove {
		delete(s.entries, key)
		s.removeFromOrder(key)
	}

	return nil
}

// Close (Method): stops the background janitor; subsequent Get/Set return ErrCacheClosed
func (s *InMemoryStore) Close() error {
	s.once.Do(func() {
		close(s.closed)
	})
	return nil
}

// Len (Method): returns the number of entries currently in the cache; exported for tests and diagnostics
func (s *InMemoryStore) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.entries)
}

// evictOldest (Method): removes the entry that was inserted first; caller must hold write lock
func (s *InMemoryStore) evictOldest() {
	if len(s.order) == 0 {
		return
	}
	oldest := s.order[0]
	s.order = s.order[1:]
	delete(s.entries, oldest)
}

// removeFromOrder (Method): removes a key from the insertion-order slice; caller must hold write lock
func (s *InMemoryStore) removeFromOrder(key string) {
	for i, k := range s.order {
		if k == key {
			s.order = append(s.order[:i], s.order[i+1:]...)
			return
		}
	}
}

// janitor (Goroutine): periodic background sweep of expired entries
func (s *InMemoryStore) janitor(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-s.closed:
			return
		case <-ticker.C:
			s.sweep()
		}
	}
}

// sweep (Method): removes all expired entries in one pass
func (s *InMemoryStore) sweep() {
	now := s.now()

	s.mu.Lock()
	defer s.mu.Unlock()

	var toRemove []string
	for key, e := range s.entries {
		if e.expired(now) {
			toRemove = append(toRemove, key)
		}
	}

	for _, key := range toRemove {
		delete(s.entries, key)
		s.removeFromOrder(key)
	}
}
