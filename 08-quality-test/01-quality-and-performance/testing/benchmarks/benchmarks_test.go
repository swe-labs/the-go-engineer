// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package benchmarks

import (
	"fmt"
	"strings"
	"testing"
)

// ============================================================================
// Stage 08: Quality and Performance - Benchmarking with testing.B
// Level: Intermediate -> Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Writing benchmarks with testing.B
//   - b.ResetTimer() for setup exclusion
//   - b.ReportAllocs() for memory profiling
//   - Sub-benchmarks for comparing approaches
//   - Running benchmarks: go test -bench=. -benchmem ./08-quality-test/01-quality-and-performance/testing/benchmarks/
//
// BENCHMARK NAMING:
//   func BenchmarkXxx(b *testing.B) - must start with Benchmark
// ============================================================================

// BenchmarkStringConcat compares three string concatenation strategies.
// This is one of the most common Go performance questions.
func BenchmarkStringConcat(b *testing.B) {
	words := make([]string, 100)
	for i := range words {
		words[i] = fmt.Sprintf("word%d", i)
	}

	// Sub-benchmark: naive += concatenation
	// This is O(n^2) because strings are immutable - each += creates a new string.
	b.Run("Plus", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var s string
			for _, w := range words {
				s += w + " "
			}
			_ = s
		}
	})

	// Sub-benchmark: strings.Join
	// O(n) - computes total length first, then copies each string once.
	b.Run("Join", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = strings.Join(words, " ")
		}
	})

	// Sub-benchmark: strings.Builder
	// O(n) - amortized, grows buffer like a dynamic array.
	// This is the idiomatic way.
	b.Run("Builder", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var builder strings.Builder
			for _, w := range words {
				builder.WriteString(w)
				builder.WriteString(" ")
			}
			_ = builder.String()
		}
	})

	// Sub-benchmark: pre-allocated Builder
	// Best performance - no reallocation.
	b.Run("BuilderPrealloc", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var builder strings.Builder
			builder.Grow(100 * 10)
			for _, w := range words {
				builder.WriteString(w)
				builder.WriteString(" ")
			}
			_ = builder.String()
		}
	})
}

// BenchmarkSliceAppend compares append with and without pre-allocation.
func BenchmarkSliceAppend(b *testing.B) {
	n := 10000

	b.Run("NoPrealloc", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var s []int
			for j := 0; j < n; j++ {
				s = append(s, j)
			}
			_ = s
		}
	})

	b.Run("WithPrealloc", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			s := make([]int, 0, n)
			for j := 0; j < n; j++ {
				s = append(s, j)
			}
			_ = s
		}
	})
}

// BenchmarkLookup compares lookup performance.
func BenchmarkLookup(b *testing.B) {
	size := 1000

	m := make(map[int]bool, size)
	s := make([]int, size)
	for i := 0; i < size; i++ {
		m[i] = true
		s[i] = i
	}
	target := size - 1

	b.Run("MapLookup", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = m[target]
		}
	})

	b.Run("SliceScan", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range s {
				if v == target {
					break
				}
			}
		}
	})
}
