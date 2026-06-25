package main

import (
	"fmt"
	"sync"
	"time"
)

type staleEntry struct {
	value      interface{}
	expiry     time.Time
	staleUntil time.Time
}

type StaleWhileRevalidateCache struct {
	mu       sync.RWMutex
	entries  map[string]staleEntry
	ttl      time.Duration
	staleTTL time.Duration
	renew    func(key string) (interface{}, error)
}

func NewStaleWhileRevalidateCache(ttl, staleTTL time.Duration, renew func(key string) (interface{}, error)) *StaleWhileRevalidateCache {
	return &StaleWhileRevalidateCache{
		entries:  make(map[string]staleEntry),
		ttl:      ttl,
		staleTTL: staleTTL,
		renew:    renew,
	}
}

func (c *StaleWhileRevalidateCache) Get(key string) (interface{}, error) {
	c.mu.RLock()
	entry, ok := c.entries[key]
	c.mu.RUnlock()

	if ok {
		if time.Now().Before(entry.expiry) {
			return entry.value, nil
		}
		if time.Now().Before(entry.staleUntil) {
			go c.refresh(key)
			return entry.value, nil
		}
	}

	return c.refresh(key)
}

func (c *StaleWhileRevalidateCache) refresh(key string) (interface{}, error) {
	value, err := c.renew(key)
	if err != nil {
		c.mu.RLock()
		entry, ok := c.entries[key]
		c.mu.RUnlock()
		if ok && time.Now().Before(entry.staleUntil) {
			return entry.value, nil
		}
		return nil, err
	}

	c.mu.Lock()
	c.entries[key] = staleEntry{
		value:      value,
		expiry:     time.Now().Add(c.ttl),
		staleUntil: time.Now().Add(c.ttl + c.staleTTL),
	}
	c.mu.Unlock()
	return value, nil
}

func (c *StaleWhileRevalidateCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

func (c *StaleWhileRevalidateCache) Invalidate(key string) {
	c.Delete(key)
}

func (c *StaleWhileRevalidateCache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}

func main() {
	renewCount := 0
	cache := NewStaleWhileRevalidateCache(
		50*time.Millisecond,
		100*time.Millisecond,
		func(key string) (interface{}, error) {
			renewCount++
			return fmt.Sprintf("value-%d", renewCount), nil
		},
	)

	val, _ := cache.Get("key")
	fmt.Println("First get:", val)

	time.Sleep(60 * time.Millisecond)

	val, _ = cache.Get("key")
	fmt.Println("Stale get (served stale, refreshing):", val)

	time.Sleep(60 * time.Millisecond)

	val, _ = cache.Get("key")
	fmt.Println("After revalidate:", val)

	fmt.Println("Total renew calls:", renewCount)
	cache.Delete("key")
	fmt.Println("After delete, size:", cache.Len())
}
