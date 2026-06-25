package main

import (
	"fmt"
	"sync"
	"testing"
)

func TestMutexCompiles(t *testing.T) {
}

func TestSafeCounter(t *testing.T) {
	c := SafeCounter{}
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Increment()
		}()
	}
	wg.Wait()
	if c.Value() != 100 {
		t.Fatalf("expected 100, got %d", c.Value())
	}
}

func TestSafeCache(t *testing.T) {
	cache := NewSafeCache()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			cache.Set(fmt.Sprintf("k-%d", n), fmt.Sprintf("v-%d", n))
		}(i)
	}
	wg.Wait()

	for i := 0; i < 100; i++ {
		v, ok := cache.Get(fmt.Sprintf("k-%d", i))
		if !ok {
			t.Fatalf("expected key k-%d to exist", i)
		}
		if v != fmt.Sprintf("v-%d", i) {
			t.Fatalf("expected v-%d, got %s", i, v)
		}
	}
}

func TestRWMutexConcurrentReads(t *testing.T) {
	cache := NewSafeCache()
	cache.Set("a", "1")

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				cache.Get("a")
			}
		}()
	}
	wg.Wait()
}
