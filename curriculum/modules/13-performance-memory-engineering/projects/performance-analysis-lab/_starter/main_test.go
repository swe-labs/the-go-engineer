package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestRateLimiter_BurstAllowanceAndThrottle(t *testing.T) {
	// 3 requests allowed per 200 milliseconds
	limiter := NewShardedRateLimiter(3, 200*time.Millisecond)

	// 1-3 should be allowed
	for i := 1; i <= 3; i++ {
		res := limiter.CheckLimit("test-client")
		if !res.Allowed {
			t.Fatalf("expected request %d to be allowed", i)
		}
		if res.Remaining != 3-i {
			t.Errorf("expected remaining %d, got %d", 3-i, res.Remaining)
		}
	}

	// 4th request must be denied
	res := limiter.CheckLimit("test-client")
	if res.Allowed {
		t.Fatalf("expected 4th request to be denied, but was allowed")
	}

	// Wait for window expiration
	time.Sleep(250 * time.Millisecond)

	// Should now be permitted again
	resAfter := limiter.CheckLimit("test-client")
	if !resAfter.Allowed {
		t.Fatalf("expected request after window reset to be allowed")
	}
}

func TestRateLimiter_DistinctClients(t *testing.T) {
	limiter := NewShardedRateLimiter(2, time.Second)

	// Exhaust client A
	limiter.CheckLimit("client-a")
	limiter.CheckLimit("client-a")
	if limiter.CheckLimit("client-a").Allowed {
		t.Fatalf("client-a should be exhausted")
	}

	// Client B should remain unaffected
	resB := limiter.CheckLimit("client-b")
	if !resB.Allowed {
		t.Fatalf("client-b should be allowed independently")
	}
}

func TestRateLimiter_ConcurrentContention(t *testing.T) {
	limiter := NewShardedRateLimiter(1000, time.Second)
	var wg sync.WaitGroup
	workers := 16
	requestsPerWorker := 500

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("tenant-%d", id%8)
			for i := 0; i < requestsPerWorker; i++ {
				_ = limiter.CheckLimit(key)
			}
		}(w)
	}

	wg.Wait()
}

func BenchmarkRateLimiter_Parallel(b *testing.B) {
	limiter := NewShardedRateLimiter(1000000, time.Minute)

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		id := 0
		for pb.Next() {
			id++
			key := fmt.Sprintf("key-%d", id%64)
			_ = limiter.CheckLimit(key)
		}
	})
}
