package main

import (
	"fmt"
	"sync"
)

type SafeCounter struct {
	mu    sync.Mutex
	value int
}

func (c *SafeCounter) Increment() {
	c.mu.Lock()
	c.value++
	c.mu.Unlock()
}

func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

type SafeCache struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewSafeCache() *SafeCache {
	return &SafeCache{data: make(map[string]string)}
}

func (c *SafeCache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.data[key]
	return v, ok
}

func (c *SafeCache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

func main() {
	fmt.Println("=== Mutex: safe counter ===")
	var wg sync.WaitGroup
	counter := SafeCounter{}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}
	wg.Wait()
	fmt.Printf("Counter value: %d (expected 100)\n", counter.Value())

	fmt.Println("\n=== RWMutex: safe cache ===")
	cache := NewSafeCache()

	var wg2 sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg2.Add(1)
		go func(n int) {
			defer wg2.Done()
			key := fmt.Sprintf("key-%d", n)
			cache.Set(key, fmt.Sprintf("val-%d", n))
		}(i)
	}
	wg2.Wait()

	var wg3 sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg3.Add(1)
		go func(n int) {
			defer wg3.Done()
			key := fmt.Sprintf("key-%d", n%10)
			if v, ok := cache.Get(key); ok {
				fmt.Printf("Read %s = %s\n", key, v)
			}
		}(i)
	}
	wg3.Wait()

	fmt.Println("\n=== Without mutex (race) ===")
	racy := 0
	var wg4 sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg4.Add(1)
		go func() {
			defer wg4.Done()
			racy++ // data race
		}()
	}
	wg4.Wait()
	fmt.Printf("Racy value: %d (non-deterministic)\n", racy)
}
