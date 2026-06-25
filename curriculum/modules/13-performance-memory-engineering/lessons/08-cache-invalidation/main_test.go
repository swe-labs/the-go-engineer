package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestCacheFreshHit(t *testing.T) {
	renew := func(key string) (interface{}, error) {
		return "fresh", nil
	}
	c := NewStaleWhileRevalidateCache(time.Minute, time.Minute, renew)

	val, err := c.Get("k")
	if err != nil {
		t.Fatal(err)
	}
	if val != "fresh" {
		t.Errorf("got %v; want fresh", val)
	}
}

func TestCacheStaleServedWithBackgroundRefresh(t *testing.T) {
	var renewCount int32
	renew := func(key string) (interface{}, error) {
		atomic.AddInt32(&renewCount, 1)
		return fmt.Sprintf("v-%d", atomic.LoadInt32(&renewCount)), nil
	}

	c := NewStaleWhileRevalidateCache(20*time.Millisecond, 100*time.Millisecond, renew)

	val1, _ := c.Get("k")
	if val1 != "v-1" {
		t.Fatalf("first get: got %v; want v-1", val1)
	}

	time.Sleep(30 * time.Millisecond)

	val2, _ := c.Get("k")
	// Should return the stale value (v-1) while triggering background refresh
	if val2 != "v-1" {
		t.Errorf("stale get: got %v; want v-1 (stale)", val2)
	}

	time.Sleep(50 * time.Millisecond)

	val3, _ := c.Get("k")
	if val3 != "v-2" {
		t.Errorf("after refresh: got %v; want v-2", val3)
	}
}

func TestCacheAfterExpiryAndStale(t *testing.T) {
	var renewCount int32
	renew := func(key string) (interface{}, error) {
		atomic.AddInt32(&renewCount, 1)
		return "new", nil
	}

	c := NewStaleWhileRevalidateCache(10*time.Millisecond, 10*time.Millisecond, renew)

	c.Get("k")
	time.Sleep(50 * time.Millisecond) // past both expiry and staleTTL

	val, err := c.Get("k")
	if err != nil {
		t.Fatal(err)
	}
	// Should have refreshed synchronously since stale is also expired
	if val != "new" {
		t.Errorf("got %v; want new", val)
	}
}

func TestCacheDelete(t *testing.T) {
	c := NewStaleWhileRevalidateCache(time.Minute, time.Minute, func(key string) (interface{}, error) {
		return "val", nil
	})

	c.Get("k")
	if c.Len() != 1 {
		t.Errorf("expected 1 entry, got %d", c.Len())
	}

	c.Delete("k")
	if c.Len() != 0 {
		t.Errorf("expected 0 after delete, got %d", c.Len())
	}

	_, err := c.Get("k")
	if err != nil {
		t.Errorf("should be able to get after delete (re-fetch): %v", err)
	}
}

func TestCacheConcurrentAccess(t *testing.T) {
	var mu sync.Mutex
	renew := func(key string) (interface{}, error) {
		mu.Lock()
		defer mu.Unlock()
		return "val", nil
	}

	c := NewStaleWhileRevalidateCache(time.Minute, time.Minute, renew)

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			key := fmt.Sprintf("k-%d", n)
			val, err := c.Get(key)
			if err != nil {
				t.Errorf("get %s: %v", key, err)
			}
			if val != "val" {
				t.Errorf("get %s: got %v; want val", key, val)
			}
		}(i)
	}
	wg.Wait()
}

func TestCacheTableDriven(t *testing.T) {
	var callCount int
	tests := []struct {
		name    string
		sleep   time.Duration
		wantVal string
		wantErr bool
	}{
		{"immediate", 0, "v1", false},
		{"stale served", 30 * time.Millisecond, "v1", false},
	}
	c := NewStaleWhileRevalidateCache(
		20*time.Millisecond,
		100*time.Millisecond,
		func(key string) (interface{}, error) {
			callCount++
			return fmt.Sprintf("v%d", callCount), nil
		},
	)

	c.Get("k")

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			time.Sleep(tc.sleep)
			val, err := c.Get("k")
			if tc.wantErr {
				if err == nil {
					t.Error("expected error")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if val != tc.wantVal {
				t.Errorf("got %v; want %v", val, tc.wantVal)
			}
		})
	}
}

func TestInvalidate(t *testing.T) {
	c := NewStaleWhileRevalidateCache(time.Minute, time.Minute, func(key string) (interface{}, error) {
		return "val", nil
	})

	c.Get("k")
	c.Invalidate("k")
	_, err := c.Get("k")
	if err != nil {
		t.Fatal(err)
	}
}
