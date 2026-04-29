// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package cache

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestGetReturnsNotFoundForMissingKey(t *testing.T) {
	t.Parallel()
	store := NewInMemoryStore(Config{MaxEntries: 8, DefaultTTL: time.Minute})
	defer store.Close()

	_, err := store.Get(context.Background(), "missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("Get error = %v, want ErrNotFound", err)
	}
}

func TestSetAndGetRoundTrip(t *testing.T) {
	t.Parallel()
	store := NewInMemoryStore(Config{MaxEntries: 8, DefaultTTL: time.Minute})
	defer store.Close()

	ctx := context.Background()
	key := TenantOrderKey(7, 101)
	value := []byte(`{"id":101,"status":"pending"}`)

	if err := store.Set(ctx, key, value, time.Minute); err != nil {
		t.Fatalf("Set returned error: %v", err)
	}

	got, err := store.Get(ctx, key)
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if string(got) != string(value) {
		t.Fatalf("Get = %q, want %q", got, value)
	}
}

func TestGetReturnsCopyNotReference(t *testing.T) {
	t.Parallel()
	store := NewInMemoryStore(Config{MaxEntries: 8, DefaultTTL: time.Minute})
	defer store.Close()

	ctx := context.Background()
	key := "test:copy"
	original := []byte("original")

	if err := store.Set(ctx, key, original, time.Minute); err != nil {
		t.Fatalf("Set returned error: %v", err)
	}

	got, _ := store.Get(ctx, key)
	got[0] = 'X' // mutate the returned value

	second, _ := store.Get(ctx, key)
	if string(second) != "original" {
		t.Fatalf("cached value was mutated: got %q, want %q", second, "original")
	}
}

func TestExpiredEntryReturnsNotFound(t *testing.T) {
	t.Parallel()
	store := NewInMemoryStore(Config{MaxEntries: 8, DefaultTTL: time.Minute})
	defer store.Close()

	// Freeze time, set entry, then advance past TTL.
	now := time.Date(2026, 4, 27, 10, 0, 0, 0, time.UTC)
	store.now = func() time.Time { return now }

	ctx := context.Background()
	key := "test:ttl"
	if err := store.Set(ctx, key, []byte("data"), 5*time.Second); err != nil {
		t.Fatalf("Set returned error: %v", err)
	}

	// Still valid.
	if _, err := store.Get(ctx, key); err != nil {
		t.Fatalf("Get before expiry returned error: %v", err)
	}

	// Advance past TTL.
	now = now.Add(6 * time.Second)
	_, err := store.Get(ctx, key)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("Get after expiry error = %v, want ErrNotFound", err)
	}

	// Lazy eviction should have removed the entry.
	if store.Len() != 0 {
		t.Fatalf("Len = %d, want 0 after lazy eviction", store.Len())
	}
}

func TestEvictsOldestWhenAtCapacity(t *testing.T) {
	t.Parallel()
	store := NewInMemoryStore(Config{MaxEntries: 2, DefaultTTL: time.Minute})
	defer store.Close()

	ctx := context.Background()
	if err := store.Set(ctx, "a", []byte("1"), time.Minute); err != nil {
		t.Fatalf("Set a: %v", err)
	}
	if err := store.Set(ctx, "b", []byte("2"), time.Minute); err != nil {
		t.Fatalf("Set b: %v", err)
	}
	// This should evict "a".
	if err := store.Set(ctx, "c", []byte("3"), time.Minute); err != nil {
		t.Fatalf("Set c: %v", err)
	}

	if _, err := store.Get(ctx, "a"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Get a after eviction error = %v, want ErrNotFound", err)
	}
	if _, err := store.Get(ctx, "b"); err != nil {
		t.Fatal("b should still be cached")
	}
	if _, err := store.Get(ctx, "c"); err != nil {
		t.Fatal("c should be cached")
	}
}

func TestDeleteRemovesEntry(t *testing.T) {
	t.Parallel()
	store := NewInMemoryStore(Config{MaxEntries: 8, DefaultTTL: time.Minute})
	defer store.Close()

	ctx := context.Background()
	key := TenantOrderKey(7, 101)
	_ = store.Set(ctx, key, []byte("data"), time.Minute)

	if err := store.Delete(ctx, key); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}

	_, err := store.Get(ctx, key)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("Get after Delete error = %v, want ErrNotFound", err)
	}
}

func TestDeletePrefixRemovesMatchingEntries(t *testing.T) {
	t.Parallel()
	store := NewInMemoryStore(Config{MaxEntries: 8, DefaultTTL: time.Minute})
	defer store.Close()

	ctx := context.Background()
	_ = store.Set(ctx, TenantOrderKey(7, 101), []byte("o1"), time.Minute)
	_ = store.Set(ctx, TenantOrderKey(7, 102), []byte("o2"), time.Minute)
	_ = store.Set(ctx, TenantPaymentListKey(7, 101), []byte("p1"), time.Minute)
	_ = store.Set(ctx, TenantOrderKey(99, 1), []byte("other"), time.Minute)

	if err := store.DeletePrefix(ctx, TenantOrderPrefix(7)); err != nil {
		t.Fatalf("DeletePrefix returned error: %v", err)
	}

	// All tenant 7 order-prefixed entries should be gone.
	for _, key := range []string{TenantOrderKey(7, 101), TenantOrderKey(7, 102), TenantPaymentListKey(7, 101)} {
		if _, err := store.Get(ctx, key); !errors.Is(err, ErrNotFound) {
			t.Fatalf("Get(%s) after prefix delete error = %v, want ErrNotFound", key, err)
		}
	}

	// Tenant 99 entry should survive.
	if _, err := store.Get(ctx, TenantOrderKey(99, 1)); err != nil {
		t.Fatalf("tenant 99 entry should survive prefix delete: %v", err)
	}
}

func TestClosedStoreRejectsOperations(t *testing.T) {
	t.Parallel()
	store := NewInMemoryStore(Config{MaxEntries: 8, DefaultTTL: time.Minute})
	store.Close()

	ctx := context.Background()
	if _, err := store.Get(ctx, "k"); !errors.Is(err, ErrCacheClosed) {
		t.Fatalf("Get on closed store error = %v, want ErrCacheClosed", err)
	}
	if err := store.Set(ctx, "k", []byte("v"), time.Minute); !errors.Is(err, ErrCacheClosed) {
		t.Fatalf("Set on closed store error = %v, want ErrCacheClosed", err)
	}
}

func TestInvalidatorInvalidatesOrderData(t *testing.T) {
	t.Parallel()
	store := NewInMemoryStore(Config{MaxEntries: 8, DefaultTTL: time.Minute})
	defer store.Close()

	ctx := context.Background()
	orderKey := TenantOrderKey(7, 101)
	listKey := TenantOrderListKey(7)
	_ = store.Set(ctx, orderKey, []byte("order"), time.Minute)
	_ = store.Set(ctx, listKey, []byte("list"), time.Minute)

	inv := NewInvalidator(store)
	inv.InvalidateOrder(ctx, 7, 101)

	if _, err := store.Get(ctx, orderKey); !errors.Is(err, ErrNotFound) {
		t.Fatal("order key should be invalidated")
	}
	if _, err := store.Get(ctx, listKey); !errors.Is(err, ErrNotFound) {
		t.Fatal("order list key should be invalidated")
	}
}

func TestInvalidatorInvalidatesPaymentData(t *testing.T) {
	t.Parallel()
	store := NewInMemoryStore(Config{MaxEntries: 8, DefaultTTL: time.Minute})
	defer store.Close()

	ctx := context.Background()
	paymentKey := TenantPaymentListKey(7, 101)
	_ = store.Set(ctx, paymentKey, []byte("payments"), time.Minute)

	inv := NewInvalidator(store)
	inv.InvalidatePayments(ctx, 7, 101)

	if _, err := store.Get(ctx, paymentKey); !errors.Is(err, ErrNotFound) {
		t.Fatal("payment key should be invalidated")
	}
}

func TestNilInvalidatorIsSafe(t *testing.T) {
	t.Parallel()
	var inv *Invalidator
	// Should not panic.
	inv.InvalidateOrder(context.Background(), 7, 101)
	inv.InvalidatePayments(context.Background(), 7, 101)
	inv.InvalidateTenantOrders(context.Background(), 7)
}

func TestNoopCacheReturnsNotFound(t *testing.T) {
	t.Parallel()
	c := NoopCache{}
	if _, err := c.Get(context.Background(), "k"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("NoopCache.Get error = %v, want ErrNotFound", err)
	}
	if err := c.Set(context.Background(), "k", []byte("v"), time.Minute); err != nil {
		t.Fatalf("NoopCache.Set returned error: %v", err)
	}
}

func TestSingleflightDeduplicatesConcurrentLoads(t *testing.T) {
	t.Parallel()

	var sf Singleflight
	var calls atomic.Int64

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			val, err := sf.Do("hot-key", func() ([]byte, error) {
				calls.Add(1)
				time.Sleep(10 * time.Millisecond) // simulate DB lookup
				return []byte("result"), nil
			})
			if err != nil {
				t.Errorf("Do returned error: %v", err)
			}
			if string(val) != "result" {
				t.Errorf("Do = %q, want %q", val, "result")
			}
		}()
	}
	wg.Wait()

	if n := calls.Load(); n != 1 {
		t.Fatalf("fn was called %d times, want 1 (singleflight should deduplicate)", n)
	}
}

func TestSingleflightPanicRecovery(t *testing.T) {
	t.Parallel()

	var sf Singleflight
	key := "panic-key"

	// First call panics - but Do should recover and return an error.
	val, err := sf.Do(key, func() ([]byte, error) {
		panic("simulated panic")
	})
	if err == nil {
		t.Fatal("expected error from panic, got nil")
	}
	if val != nil {
		t.Fatalf("expected nil value from panic, got %q", val)
	}
	if got := err.Error(); got != "cache: singleflight panic: simulated panic" {
		t.Fatalf("unexpected error message: %s", got)
	}

	// Second call should not block, meaning the key was cleaned up.
	done := make(chan struct{})
	go func() {
		v, e := sf.Do(key, func() ([]byte, error) {
			return []byte("recovered"), nil
		})
		if e != nil {
			t.Errorf("Do returned error: %v", e)
		}
		if string(v) != "recovered" {
			t.Errorf("Do = %q, want %q", v, "recovered")
		}
		close(done)
	}()

	select {
	case <-done:
		// Success
	case <-time.After(1 * time.Second):
		t.Fatal("second Do call blocked forever, panic cleanup failed")
	}
}

func TestEmptyKeyReturnsError(t *testing.T) {
	t.Parallel()
	store := NewInMemoryStore(Config{MaxEntries: 8, DefaultTTL: time.Minute})
	defer store.Close()

	ctx := context.Background()
	if _, err := store.Get(ctx, ""); !errors.Is(err, ErrKeyEmpty) {
		t.Fatalf("Get empty key error = %v, want ErrKeyEmpty", err)
	}
	if err := store.Set(ctx, "", []byte("v"), time.Minute); !errors.Is(err, ErrKeyEmpty) {
		t.Fatalf("Set empty key error = %v, want ErrKeyEmpty", err)
	}
}

func TestUpdateExistingKeyDoesNotChangeCapacity(t *testing.T) {
	t.Parallel()
	store := NewInMemoryStore(Config{MaxEntries: 2, DefaultTTL: time.Minute})
	defer store.Close()

	ctx := context.Background()
	_ = store.Set(ctx, "a", []byte("1"), time.Minute)
	_ = store.Set(ctx, "b", []byte("2"), time.Minute)

	// Update "a" - should NOT evict "b".
	_ = store.Set(ctx, "a", []byte("updated"), time.Minute)

	if store.Len() != 2 {
		t.Fatalf("Len = %d, want 2", store.Len())
	}
	got, _ := store.Get(ctx, "a")
	if string(got) != "updated" {
		t.Fatalf("Get a = %q, want %q", got, "updated")
	}
	if _, err := store.Get(ctx, "b"); err != nil {
		t.Fatal("b should still be cached after updating a")
	}
}
