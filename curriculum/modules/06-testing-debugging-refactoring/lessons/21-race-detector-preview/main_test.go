package main

import (
	"sync"
	"testing"
)

func TestCounterSafe(t *testing.T) {
	c := &Counter{}
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.mu.Lock()
			c.value++
			c.mu.Unlock()
		}()
	}
	wg.Wait()
	if c.value != 10 {
		t.Errorf("expected 10, got %d", c.value)
	}
}

func TestCounterParallel(t *testing.T) {
	c := &Counter{}
	var wg sync.WaitGroup
	n := 50
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.mu.Lock()
			c.value++
			c.mu.Unlock()
		}()
	}
	wg.Wait()
	if c.value != n {
		t.Errorf("expected %d, got %d", n, c.value)
	}
}

func TestUnsafeIncrementRaces(t *testing.T) {
	c := &Counter{}
	// This intentionally has a race — run with -race to detect it.
	UnsafeIncrement(c)
	// We don't assert on the value because it's non-deterministic.
	t.Logf("unsafe counter value: %d (non-deterministic)", c.value)
}
