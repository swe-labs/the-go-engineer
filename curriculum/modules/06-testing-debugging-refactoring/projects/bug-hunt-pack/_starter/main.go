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
// BUG: SafeCounter is not thread-safe. Concurrent calls to Inc() produce data races.
// TODO: Fix by protecting count with a sync.Mutex.

type SafeCounter struct {
	count int
}

func NewSafeCounter() *SafeCounter {
	return &SafeCounter{}
}

func (c *SafeCounter) Inc() {
	// BUG: Unsynchronized mutation
	c.count++
}

func (c *SafeCounter) Value() int {
	return c.count
}

// --- Bug 2: Safe Memory Cache ---
// BUG: If KeyValueCache is created as a zero-value struct (var cache KeyValueCache),
// calling Set() panics with "assignment to entry in nil map".
// TODO: Guard or lazily initialize c.data before writing.

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
	// BUG: Panics if c.data is nil
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

// --- Bug 3: Safe Pagination Slicer ---
// BUG: If offset + limit exceeds len(items), slicing items[offset : offset+limit] panics
// with "runtime error: slice bounds out of range".
// TODO: Clamp offset and end bounds safely.

func PaginateSlice[T any](items []T, offset, limit int) []T {
	if offset < 0 {
		offset = 0
	}
	if offset >= len(items) {
		return []T{}
	}

	// BUG: Does not check if offset+limit > len(items)
	return items[offset : offset+limit]
}

// --- Bug 4: Shadowed Error Variable ---
// BUG: The short variable declaration (:=) on line formatName shadows the outer `err`,
// causing the function to return a nil error even when formatting fails!
// TODO: Fix variable shadowing so the error returned by formatName is returned to the caller.

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

	// BUG: err is shadowed here!
	formatted, err := formatName(name)
	if err != nil {
		// Supposed to return here, but if changed or falling through, outer err is nil
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
}
