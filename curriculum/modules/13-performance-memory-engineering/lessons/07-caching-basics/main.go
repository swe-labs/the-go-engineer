package main

import (
	"fmt"
	"sync"
	"time"
)

type CacheEntry struct {
	value  interface{}
	expiry time.Time
}

type TTLCache struct {
	mu       sync.RWMutex
	entries  map[string]CacheEntry
	ttl      time.Duration
	stopChan chan struct{}
}

func NewTTLCache(ttl time.Duration, cleanupInterval time.Duration) *TTLCache {
	c := &TTLCache{
		entries:  make(map[string]CacheEntry),
		ttl:      ttl,
		stopChan: make(chan struct{}),
	}
	go c.cleanupLoop(cleanupInterval)
	return c
}

func (c *TTLCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = CacheEntry{
		value:  value,
		expiry: time.Now().Add(c.ttl),
	}
}

func (c *TTLCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	entry, ok := c.entries[key]
	c.mu.RUnlock()
	if !ok {
		return nil, false
	}
	if time.Now().After(entry.expiry) {
		c.mu.Lock()
		delete(c.entries, key)
		c.mu.Unlock()
		return nil, false
	}
	return entry.value, true
}

func (c *TTLCache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}

func (c *TTLCache) cleanupLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.deleteExpired()
		case <-c.stopChan:
			return
		}
	}
}

func (c *TTLCache) deleteExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	for k, entry := range c.entries {
		if now.After(entry.expiry) {
			delete(c.entries, k)
		}
	}
}

func (c *TTLCache) Stop() {
	close(c.stopChan)
}

func main() {
	cache := NewTTLCache(100*time.Millisecond, 50*time.Millisecond)
	defer cache.Stop()

	cache.Set("name", "Alice")
	if val, ok := cache.Get("name"); ok {
		fmt.Println("Got:", val)
	} else {
		fmt.Println("Miss")
	}

	time.Sleep(200 * time.Millisecond)

	if _, ok := cache.Get("name"); !ok {
		fmt.Println("Expired after TTL")
	}

	cache.Set("count", 42)
	fmt.Println("Cache size:", cache.Len())
}
