package main

import "testing"

func BenchmarkPR5Summary(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = pr_5Summary("  benchmark input  ")
	}
}
