// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import "testing"

// ============================================================================
// Benchmarks — run with: go test -bench=. -benchmem ./24-errgroup-and-pools/3-sync-pool
// ============================================================================

func BenchmarkWithPool(b *testing.B) {
	body := `{"status":"ok","data":{"user_id":42}}`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = buildHTTPResponseWithPool(200, body)
	}
}

func BenchmarkWithoutPool(b *testing.B) {
	body := `{"status":"ok","data":{"user_id":42}}`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = buildHTTPResponseWithoutPool(200, body)
	}
}
