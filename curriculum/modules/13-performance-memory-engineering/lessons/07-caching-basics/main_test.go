package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestCacheSetAndGet(t *testing.T) {
	c := NewTTLCache(time.Minute, time.Minute)
	defer c.Stop()

	c.Set("key1", "value1")
	val, ok := c.Get("key1")
	if !ok {
		t.Fatal("expected to find key1")
	}
	if val != "value1" {
		t.Errorf("got %v; want value1", val)
	}
}

func TestCacheMiss(t *testing.T) {
	c := NewTTLCache(time.Minute, time.Minute)
	defer c.Stop()

	_, ok := c.Get("nonexistent")
	if ok {
		t.Error("expected miss for nonexistent key")
	}
}

func TestCacheTTLExpiry(t *testing.T) {
	c := NewTTLCache(50*time.Millisecond, 10*time.Millisecond)
	defer c.Stop()

	c.Set("key", "value")
	time.Sleep(100 * time.Millisecond)

	_, ok := c.Get("key")
	if ok {
		t.Error("expected key to expire")
	}
}

func TestCacheConcurrentAccess(t *testing.T) {
	c := NewTTLCache(time.Minute, time.Minute)
	defer c.Stop()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			key := fmt.Sprintf("key-%d", n)
			c.Set(key, n)
			val, ok := c.Get(key)
			if !ok {
				t.Errorf("miss for %s", key)
			} else if val != n {
				t.Errorf("got %v; want %d", val, n)
			}
		}(i)
	}
	wg.Wait()

	if c.Len() != 100 {
		t.Errorf("expected 100 entries, got %d", c.Len())
	}
}

func TestCacheLen(t *testing.T) {
	c := NewTTLCache(time.Minute, time.Minute)
	defer c.Stop()

	if c.Len() != 0 {
		t.Errorf("expected 0, got %d", c.Len())
	}
	c.Set("a", 1)
	c.Set("b", 2)
	if c.Len() != 2 {
		t.Errorf("expected 2, got %d", c.Len())
	}
}

func TestCacheOverwrite(t *testing.T) {
	c := NewTTLCache(time.Minute, time.Minute)
	defer c.Stop()

	c.Set("key", "old")
	c.Set("key", "new")
	val, ok := c.Get("key")
	if !ok {
		t.Fatal("expected to find key")
	}
	if val != "new" {
		t.Errorf("got %v; want new", val)
	}
}

func TestCacheTableDriven(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		value   interface{}
		wantOK  bool
		wantVal interface{}
	}{
		{"string value", "s", "hello", true, "hello"},
		{"int value", "i", 42, true, 42},
		{"nil value", "n", nil, true, nil},
	}
	c := NewTTLCache(time.Minute, time.Minute)
	defer c.Stop()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c.Set(tc.key, tc.value)
			val, ok := c.Get(tc.key)
			if ok != tc.wantOK {
				t.Errorf("ok = %v; want %v", ok, tc.wantOK)
			}
			if val != tc.wantVal {
				t.Errorf("val = %v; want %v", val, tc.wantVal)
			}
		})
	}
}

func TestCleanupRemovesExpired(t *testing.T) {
	c := NewTTLCache(20*time.Millisecond, 10*time.Millisecond)
	defer c.Stop()

	c.Set("a", 1)
	c.Set("b", 2)

	time.Sleep(50 * time.Millisecond)

	if c.Len() != 0 {
		t.Errorf("expected 0 after cleanup, got %d", c.Len())
	}
}

func TestGetExtendsNothing(t *testing.T) {
	c := NewTTLCache(50*time.Millisecond, 100*time.Millisecond)
	defer c.Stop()

	c.Set("key", "val")
	time.Sleep(30 * time.Millisecond)

	_, ok := c.Get("key")
	if !ok {
		t.Fatal("expected hit before TTL")
	}

	time.Sleep(40 * time.Millisecond)

	_, ok = c.Get("key")
	if ok {
		t.Error("expected miss after TTL; Get should not extend TTL")
	}
}
