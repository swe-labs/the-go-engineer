package main

import (
	"fmt"
	"sync"
	"time"
)

type CacheEntry struct {
	Value     string
	ExpiresAt time.Time
}

type CacheAside struct {
	mu       sync.RWMutex
	entries  map[string]CacheEntry
	ttl      time.Duration
	loadFunc func(key string) (string, error)
	misses   int
	hits     int
}

func NewCacheAside(ttl time.Duration, load func(string) (string, error)) *CacheAside {
	return &CacheAside{
		entries:  make(map[string]CacheEntry),
		ttl:      ttl,
		loadFunc: load,
	}
}

func (c *CacheAside) Get(key string) (string, error) {
	c.mu.RLock()
	entry, ok := c.entries[key]
	c.mu.RUnlock()

	if ok && time.Now().Before(entry.ExpiresAt) {
		c.mu.Lock()
		c.hits++
		c.mu.Unlock()
		return entry.Value, nil
	}

	c.mu.Lock()
	c.misses++
	c.mu.Unlock()

	value, err := c.loadFunc(key)
	if err != nil {
		return "", err
	}
	c.mu.Lock()
	c.entries[key] = CacheEntry{Value: value, ExpiresAt: time.Now().Add(c.ttl)}
	c.mu.Unlock()
	return value, nil
}

func (c *CacheAside) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = CacheEntry{Value: value, ExpiresAt: time.Now().Add(c.ttl)}
}

func (c *CacheAside) Invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

func (c *CacheAside) Stats() (hits, misses int) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.hits, c.misses
}

func (c *CacheAside) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}

type WriteBehindCache struct {
	mu        sync.Mutex
	entries   map[string]string
	writeFunc func(key, value string) error
	pending   int
}

func NewWriteBehindCache(write func(string, string) error) *WriteBehindCache {
	return &WriteBehindCache{
		entries:   make(map[string]string),
		writeFunc: write,
	}
}

func (w *WriteBehindCache) Set(key, value string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.entries[key] = value
	w.pending++
}

func (w *WriteBehindCache) Get(key string) (string, bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	v, ok := w.entries[key]
	return v, ok
}

func (w *WriteBehindCache) Flush() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	for k, v := range w.entries {
		if err := w.writeFunc(k, v); err != nil {
			return err
		}
		delete(w.entries, k)
	}
	w.pending = 0
	return nil
}

func (w *WriteBehindCache) Pending() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.pending
}

type CacheHierarchy struct {
	l1 *CacheAside
	l2 *CacheAside
}

func NewCacheHierarchy(l1ttl, l2ttl time.Duration, load func(string) (string, error)) *CacheHierarchy {
	return &CacheHierarchy{
		l1: NewCacheAside(l1ttl, func(key string) (string, error) {
			val, err := load(key)
			return val, err
		}),
		l2: NewCacheAside(l2ttl, load),
	}
}

func (h *CacheHierarchy) Get(key string) (string, error) {
	val, err := h.l1.Get(key)
	if err == nil {
		return val, nil
	}
	val, err = h.l2.Get(key)
	if err != nil {
		return "", err
	}
	h.l1.Set(key, val)
	return val, nil
}

func expensiveLoad(key string) (string, error) {
	time.Sleep(10 * time.Millisecond)
	return fmt.Sprintf("value-for-%s", key), nil
}

func main() {
	cache := NewCacheAside(100*time.Millisecond, expensiveLoad)

	val, _ := cache.Get("user:1")
	fmt.Println("first get:", val)

	val, _ = cache.Get("user:1")
	fmt.Println("second get (cache hit):", val)

	hits, misses := cache.Stats()
	fmt.Printf("hits=%d misses=%d\n", hits, misses)

	time.Sleep(150 * time.Millisecond)
	val, _ = cache.Get("user:1")
	fmt.Println("after expiry:", val)

	wb := NewWriteBehindCache(func(k, v string) error {
		fmt.Printf("flushed: %s=%s\n", k, v)
		return nil
	})
	wb.Set("k1", "v1")
	wb.Set("k2", "v2")
	fmt.Printf("pending before flush: %d\n", wb.Pending())
	wb.Flush()
	fmt.Printf("pending after flush: %d\n", wb.Pending())

	_ = CacheHierarchy{}
}
