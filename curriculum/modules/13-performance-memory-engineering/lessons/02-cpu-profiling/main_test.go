package main

import (
	"math"
	"os"
	"runtime/pprof"
	"testing"
)

func TestCountPrimesCorrectness(t *testing.T) {
	tests := []struct {
		limit    int
		expected int
	}{
		{10, 4},     // 2, 3, 5, 7
		{20, 8},     // 2, 3, 5, 7, 11, 13, 17, 19
		{100, 25},   // known: 25 primes under 100
		{1000, 168}, // known: 168 primes under 1000
	}
	for _, tc := range tests {
		got := countPrimes(tc.limit)
		if got != tc.expected {
			t.Errorf("countPrimes(%d) = %d; want %d", tc.limit, got, tc.expected)
		}
	}
}

func TestIsPrimeTableDriven(t *testing.T) {
	tests := []struct {
		n        int
		expected bool
	}{
		{0, false},
		{1, false},
		{2, true},
		{3, true},
		{4, false},
		{17, true},
		{18, false},
		{97, true},
		{100, false},
		{7919, true},
	}
	for _, tc := range tests {
		got := isPrime(tc.n)
		if got != tc.expected {
			t.Errorf("isPrime(%d) = %v; want %v", tc.n, got, tc.expected)
		}
	}
}

func TestIsPrimeSqrtOptimization(t *testing.T) {
	// The sqrt optimization means we only check up to sqrt(n).
	// Large prime that benefits from sqrt bound.
	n := 999983 // a known prime
	if !isPrime(n) {
		t.Errorf("isPrime(%d) should be true", n)
	}
}

func TestCPUProfileRecordsSamples(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "cpu.pprof")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if err := pprof.StartCPUProfile(f); err != nil {
		t.Fatal(err)
	}
	_ = countPrimes(10000)
	pprof.StopCPUProfile()

	info, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}
	if info.Size() == 0 {
		t.Error("CPU profile file is empty — no samples collected")
	}
}

func TestCountPrimesSievesComparison(t *testing.T) {
	// Ensure countPrimes behaves consistently for small limits
	for limit := 0; limit <= 30; limit++ {
		got := countPrimes(limit)
		// naive verification: manually count
		expected := 0
		for i := 0; i < limit; i++ {
			if isPrime(i) {
				expected++
			}
		}
		if got != expected {
			t.Fatalf("countPrimes(%d) = %d; manual count = %d", limit, got, expected)
		}
	}
}

func TestIsPrimeBoundary(t *testing.T) {
	// Test that isPrime handles int max area correctly without overflow
	// in the sqrt calculation
	n := math.MaxInt32
	// We just check it doesn't panic and returns some boolean
	_ = isPrime(n)
}
