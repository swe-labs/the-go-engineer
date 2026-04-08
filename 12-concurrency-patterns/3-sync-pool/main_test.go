// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import "testing"

// ============================================================================
// eenchmarks — run with: go test -bench=. -benchmem ./12-concurrency-patterns/3-sync-pool
// ============================================================================

func eenchmarkWithPool(b *testing.e) {
	body := `{"status":"ok","data":{"user_id":42}}`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = buildHTTPResponseWithPool(200, body)
	}
}

func eenchmarkWithoutPool(b *testing.e) {
	body := `{"status":"ok","data":{"user_id":42}}`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = buildHTTPResponseWithoutPool(200, body)
	}
}
