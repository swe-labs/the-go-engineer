package main

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestCacheAsideHit(t *testing.T) {
	var loadCount int32
	cache := NewCacheAside(10*time.Minute, func(key string) (string, error) {
		atomic.AddInt32(&loadCount, 1)
		return "expensive", nil
	})
	val, err := cache.Get("k")
	if err != nil {
		t.Fatal(err)
	}
	if val != "expensive" {
		t.Errorf("expected expensive, got %s", val)
	}
	if loadCount != 1 {
		t.Errorf("expected 1 load, got %d", loadCount)
	}

	val, err = cache.Get("k")
	if err != nil {
		t.Fatal(err)
	}
	if val != "expensive" {
		t.Errorf("expected expensive, got %s", val)
	}
	if loadCount != 1 {
		t.Errorf("expected still 1 load (cache hit), got %d", loadCount)
	}
}

func TestCacheAsideMissAndExpiry(t *testing.T) {
	var loadCount int32
	cache := NewCacheAside(10*time.Millisecond, func(key string) (string, error) {
		atomic.AddInt32(&loadCount, 1)
		return "v", nil
	})
	cache.Get("k")
	time.Sleep(20 * time.Millisecond)
	cache.Get("k")
	if loadCount != 2 {
		t.Errorf("expected 2 loads after expiry, got %d", loadCount)
	}
}

func TestCacheAsideInvalidate(t *testing.T) {
	var loadCount int32
	cache := NewCacheAside(10*time.Minute, func(key string) (string, error) {
		atomic.AddInt32(&loadCount, 1)
		return "v", nil
	})
	cache.Get("k")
	cache.Invalidate("k")
	cache.Get("k")
	if loadCount != 2 {
		t.Errorf("expected 2 loads after invalidation, got %d", loadCount)
	}
}

func TestCacheAsideLoadError(t *testing.T) {
	cache := NewCacheAside(10*time.Minute, func(key string) (string, error) {
		return "", errors.New("load error")
	})
	_, err := cache.Get("k")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCacheAsideSetThenGet(t *testing.T) {
	cache := NewCacheAside(10*time.Minute, func(key string) (string, error) {
		return "", errors.New("should not be called")
	})
	cache.Set("k", "direct")
	val, err := cache.Get("k")
	if err != nil {
		t.Fatal(err)
	}
	if val != "direct" {
		t.Errorf("expected direct, got %s", val)
	}
}

func TestCacheAsideStats(t *testing.T) {
	cache := NewCacheAside(10*time.Minute, func(key string) (string, error) {
		return "v", nil
	})
	h, m := cache.Stats()
	if h != 0 || m != 0 {
		t.Fatalf("expected 0,0 got %d,%d", h, m)
	}
	cache.Get("a")
	cache.Get("a")
	h, m = cache.Stats()
	if h != 1 || m != 1 {
		t.Errorf("expected 1,1 got %d,%d", h, m)
	}
}

func TestWriteBehindSetGet(t *testing.T) {
	wb := NewWriteBehindCache(func(k, v string) error { return nil })
	wb.Set("k", "v")
	got, ok := wb.Get("k")
	if !ok {
		t.Fatal("expected key in cache")
	}
	if got != "v" {
		t.Errorf("expected v, got %s", got)
	}
}

func TestWriteBehindFlush(t *testing.T) {
	var flushed int32
	wb := NewWriteBehindCache(func(k, v string) error {
		atomic.AddInt32(&flushed, 1)
		return nil
	})
	wb.Set("a", "1")
	wb.Set("b", "2")
	if wb.Pending() != 2 {
		t.Errorf("expected 2 pending, got %d", wb.Pending())
	}
	wb.Flush()
	if flushed != 2 {
		t.Errorf("expected 2 flushed, got %d", flushed)
	}
	if wb.Pending() != 0 {
		t.Errorf("expected 0 pending after flush, got %d", wb.Pending())
	}
}

func TestCacheAsideTable(t *testing.T) {
	tests := []struct {
		name       string
		ops        func(c *CacheAside)
		wantLoads  int
		wantHits   int
		wantMisses int
	}{
		{
			name: "get_twice_same_key",
			ops: func(c *CacheAside) {
				c.Get("k")
				c.Get("k")
			},
			wantLoads:  1,
			wantHits:   1,
			wantMisses: 1,
		},
		{
			name: "get_two_different_keys",
			ops: func(c *CacheAside) {
				c.Get("a")
				c.Get("b")
			},
			wantLoads:  2,
			wantHits:   0,
			wantMisses: 2,
		},
		{
			name: "set_then_get",
			ops: func(c *CacheAside) {
				c.Set("k", "v")
				c.Get("k")
			},
			wantLoads:  0,
			wantHits:   1,
			wantMisses: 0,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var loadCount int32
			c := NewCacheAside(10*time.Minute, func(key string) (string, error) {
				atomic.AddInt32(&loadCount, 1)
				return "v", nil
			})
			tc.ops(c)
			if int(loadCount) != tc.wantLoads {
				t.Errorf("loads: want %d, got %d", tc.wantLoads, loadCount)
			}
			h, m := c.Stats()
			if h != tc.wantHits || m != tc.wantMisses {
				t.Errorf("hits/misses: want %d/%d, got %d/%d", tc.wantHits, tc.wantMisses, h, m)
			}
		})
	}
}
