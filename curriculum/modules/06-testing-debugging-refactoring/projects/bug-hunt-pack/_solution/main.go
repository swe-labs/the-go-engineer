package main

import (
	"errors"
	"fmt"
	"sync"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
)

// --- Bug 1: Concurrency Safe Accumulator ---
// Fixed: Thread-safe accumulator with sync.Mutex preventing data races.

type SafeCounter struct {
	mu    sync.Mutex
	count int
}

func NewSafeCounter() *SafeCounter {
	return &SafeCounter{}
}

func (c *SafeCounter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

// --- Bug 2: Safe Memory Cache (Fixing Nil Map Assignment) ---
// Fixed: Cache ensures underlying map is initialized before insertion.

type KeyValueCache struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewKeyValueCache() *KeyValueCache {
	return &KeyValueCache{
		data: make(map[string]string),
	}
}

func (c *KeyValueCache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.data == nil {
		c.data = make(map[string]string)
	}
	c.data[key] = value
}

func (c *KeyValueCache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.data == nil {
		return "", false
	}
	val, ok := c.data[key]
	return val, ok
}

// --- Bug 3: Safe Pagination Slicer (Fixing Slice Bounds Out of Range) ---
// Fixed: Off-by-one and slice bounds clamped to len(items).

func PaginateSlice[T any](items []T, offset, limit int) []T {
	if offset < 0 {
		offset = 0
	}
	if offset >= len(items) {
		return []T{}
	}
	if limit <= 0 {
		return []T{}
	}

	end := offset + limit
	if end > len(items) {
		end = len(items)
	}

	return items[offset:end]
}

// --- Bug 4: Shadowed Error Fix ---
// Fixed: Error variable is not shadowed by short variable declaration (:=).

type User struct {
	ID   string
	Name string
}

func ValidateAndFormatUser(id, name string) (User, error) {
	var err error

	if id == "" {
		err = fmt.Errorf("%w: user id cannot be empty", ErrInvalidInput)
		return User{}, err
	}

	// Previously: formatted, err := formatName(name) (shadowed outer err!)
	var formatted string
	formatted, err = formatName(name)
	if err != nil {
		return User{}, err
	}

	return User{ID: id, Name: formatted}, nil
}

func formatName(name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("%w: name cannot be empty", ErrInvalidInput)
	}
	return "User: " + name, nil
}

func main() {
	counter := NewSafeCounter()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Inc()
		}()
	}
	wg.Wait()
	fmt.Printf("SafeCounter: %d\n", counter.Value())

	cache := NewKeyValueCache()
	cache.Set("status", "operational")
	val, _ := cache.Get("status")
	fmt.Printf("Cache status: %s\n", val)

	items := []int{10, 20, 30, 40, 50}
	page := PaginateSlice(items, 2, 10) // requesting 10 items starting at 2
	fmt.Printf("Paginated: %v\n", page)

	user, err := ValidateAndFormatUser("u100", "Grace Hopper")
	if err == nil {
		fmt.Printf("Validated: %+v\n", user)
	}
}
