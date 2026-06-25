package main

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestAtomicsCompiles(t *testing.T) {
}

func TestAtomicCounter(t *testing.T) {
	var wg sync.WaitGroup
	counter := AtomicCounter{}
	n := 1000
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Add(1)
		}()
	}
	wg.Wait()
	if counter.Load() != int64(n) {
		t.Fatalf("expected %d, got %d", n, counter.Load())
	}
}

func TestAtomicFlagCAS(t *testing.T) {
	flag := AtomicFlag{}
	if flag.IsSet() {
		t.Fatal("expected false initially")
	}
	if !flag.TrySet() {
		t.Fatal("expected TrySet to succeed")
	}
	if flag.TrySet() {
		t.Fatal("expected TrySet to fail on second call")
	}
}

func TestAtomicFlagConcurrent(t *testing.T) {
	flag := AtomicFlag{}
	var wg sync.WaitGroup
	var successes int64

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if flag.TrySet() {
				atomic.AddInt64(&successes, 1)
			}
		}()
	}
	wg.Wait()
	if successes != 1 {
		t.Fatalf("expected exactly 1 success, got %d", successes)
	}
}

func TestAtomicLoadStore(t *testing.T) {
	counter := AtomicCounter{}
	counter.Add(42)
	if counter.Load() != 42 {
		t.Fatalf("expected 42, got %d", counter.Load())
	}
}
