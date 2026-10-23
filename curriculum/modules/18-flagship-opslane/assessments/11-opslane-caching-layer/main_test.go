package main

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestCachedOrderService_ReadThroughAndHits(t *testing.T) {
	db := NewMockDatabaseStore()
	db.orders["ord-hit-1"] = &OrderRecord{ID: "ord-hit-1", Title: "Kubernetes Cluster Provisioning", Amount: 2400.0}

	cache := NewTTLCache()
	service := NewCachedOrderService(db, cache, 5*time.Minute)

	ctx := context.Background()

	// 1. Initial read -> Cache Miss -> DB query
	o1, err := service.GetOrder(ctx, "ord-hit-1")
	if err != nil || o1 == nil {
		t.Fatalf("expected order to be fetched, got: %v, err: %v", o1, err)
	}
	if service.dbReads.Load() != 1 {
		t.Errorf("expected 1 DB read on initial miss, got %d", service.dbReads.Load())
	}

	// 2. Second read -> Cache Hit -> DB NOT queried
	o2, err := service.GetOrder(ctx, "ord-hit-1")
	if err != nil || o2 == nil {
		t.Fatalf("expected order to be returned, got: %v", err)
	}
	if service.dbReads.Load() != 1 {
		t.Errorf("expected DB reads to remain 1 on cache hit, got %d", service.dbReads.Load())
	}
}

func TestCachedOrderService_InvalidationOnUpdate(t *testing.T) {
	// Core practice task requirement:
	// Read a cached resource, modify it through an update operation,
	// and verify that the subsequent read returns the fresh value rather than the stale cached value.
	db := NewMockDatabaseStore()
	db.orders["ord-inv-2"] = &OrderRecord{ID: "ord-inv-2", Title: "Old Version 1.0", Amount: 100.0}

	cache := NewTTLCache()
	service := NewCachedOrderService(db, cache, 5*time.Minute)

	ctx := context.Background()

	// 1. Warm cache
	firstRead, err := service.GetOrder(ctx, "ord-inv-2")
	if err != nil || firstRead.Title != "Old Version 1.0" {
		t.Fatalf("expected initial title 'Old Version 1.0', got: %+v", firstRead)
	}

	// 2. Perform Update mutation through service
	updatedRecord := &OrderRecord{
		ID:     "ord-inv-2",
		Title:  "Fresh Version 2.0 (Updated)",
		Amount: 150.0,
	}
	if err := service.UpdateOrder(ctx, updatedRecord); err != nil {
		t.Fatalf("failed to update order: %v", err)
	}

	// 3. Subsequent read MUST return the fresh value (cache invalidation verified)
	subsequentRead, err := service.GetOrder(ctx, "ord-inv-2")
	if err != nil {
		t.Fatalf("failed to read order after update: %v", err)
	}

	if subsequentRead.Title != "Fresh Version 2.0 (Updated)" || subsequentRead.Amount != 150.0 {
		t.Fatalf("STALE CACHE DEFECT: expected fresh title 'Fresh Version 2.0 (Updated)', but received stale value: %+v", subsequentRead)
	}

	// Database read count should now be 2 (initial fetch + post-invalidation fetch)
	if service.dbReads.Load() != 2 {
		t.Errorf("expected 2 DB reads, got %d", service.dbReads.Load())
	}
}

func TestCachedOrderService_SingleflightStampede(t *testing.T) {
	db := NewMockDatabaseStore()
	db.orders["ord-stampede-3"] = &OrderRecord{ID: "ord-stampede-3", Title: "Flash Sale Special", Amount: 19.99}

	cache := NewTTLCache()
	service := NewCachedOrderService(db, cache, 1*time.Minute)

	const concurrentCallers = 25
	var wg sync.WaitGroup
	wg.Add(concurrentCallers)

	start := make(chan struct{})
	for i := 0; i < concurrentCallers; i++ {
		go func() {
			defer wg.Done()
			<-start
			_, _ = service.GetOrder(context.Background(), "ord-stampede-3")
		}()
	}
	close(start)

	wg.Wait()

	// Singleflight must coalesce all concurrent initial reads into a single DB query
	reads := service.dbReads.Load()
	if reads != 1 {
		t.Errorf("expected singleflight to execute exactly 1 DB query, but observed %d", reads)
	}
}

func TestTTLCache_Expiration(t *testing.T) {
	cache := NewTTLCache()
	cache.Set("short-lived-key", []byte("payload"), 10*time.Millisecond)

	// Immediate lookup -> Hit
	if _, hit := cache.Get("short-lived-key"); !hit {
		t.Fatalf("expected immediate cache hit")
	}

	// Wait past TTL
	time.Sleep(25 * time.Millisecond)

	// Expired lookup -> Miss
	if _, hit := cache.Get("short-lived-key"); hit {
		t.Errorf("expected cache entry to expire after TTL")
	}
}
