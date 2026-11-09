package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type AtomicCounter struct {
	value atomic.Int64
}

func (c *AtomicCounter) Add(n int64) {
	c.value.Add(n)
}

func (c *AtomicCounter) Load() int64 {
	return c.value.Load()
}

type Config struct {
	mu     sync.RWMutex
	params map[string]string
}

func (c *Config) Get(key string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.params[key]
}

func (c *Config) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.params[key] = value
}

type AtomicFlag struct {
	value atomic.Bool
}

func (f *AtomicFlag) Set() {
	f.value.Store(true)
}

func (f *AtomicFlag) IsSet() bool {
	return f.value.Load()
}

func (f *AtomicFlag) TrySet() bool {
	return f.value.CompareAndSwap(false, true)
}

func main() {
	fmt.Println("=== Atomic counter ===")
	var wg sync.WaitGroup
	counter := AtomicCounter{}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Add(1)
		}()
	}
	wg.Wait()
	fmt.Printf("Counter: %d (expected 1000)\n", counter.Load())

	fmt.Println("\n=== Atomic flag (CompareAndSwap) ===")
	flag := AtomicFlag{}

	var wg2 sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg2.Add(1)
		go func(id int) {
			defer wg2.Done()
			if flag.TrySet() {
				fmt.Printf("goroutine %d set the flag\n", id)
			} else {
				fmt.Printf("goroutine %d: flag already set\n", id)
			}
		}(i)
	}
	wg2.Wait()

	fmt.Println("\n=== When to use mutex vs atomic ===")
	config := Config{params: make(map[string]string)}
	config.Set("host", "localhost")
	config.Set("port", "8080")

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = config.Get("host")
		}()
	}
	wg.Wait()
	fmt.Printf("Config host: %s\n", config.Get("host"))
}
