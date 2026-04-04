// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime"
	"runtime/pprof"
	"strings"
)

// ============================================================================
// Section 25: Profiling — CPU Profile
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - runtime/pprof: writing CPU and memory profiles to files
//   - go tool pprof: reading profiles in the terminal and as a web UI
//   - Identifying hot functions and optimising them
//   - The before/after pattern: benchmark → profile → fix → verify
//
// AFTER RUNNING THIS PROGRAM:
//   1. Open the CPU profile:
//      go tool pprof cpu.prof
//      (pprof) top10          -- show top 10 functions by CPU time
//      (pprof) list processLogs -- show annotated source for processLogs
//      (pprof) web            -- open flame graph (requires graphviz)
//
//   2. Open as a web UI (much better):
//      go tool pprof -http=:8090 cpu.prof
//      Browse to http://localhost:8090
//
// ENGINEERING DEPTH:
//   CPU profiling works by installing a signal handler (SIGPROF) that fires
//   100 times per second. Each time it fires, the Go runtime records the current
//   goroutine's call stack. After profiling ends, the tool counts how often each
//   function appeared in a stack sample. High sample count = lots of CPU time.
//
//   "flat" time: time spent executing this function (not its callees)
//   "cum" (cumulative) time: time in this function + everything it calls
//
// RUN: go run ./25-profiling/1-cpu-profile
// ============================================================================

// ============================================================================
// Intentionally slow implementations — we will profile and fix them
// ============================================================================

// slowLogProcessor compiles a regex and concatenates strings inside a hot loop.
// Both are classic Go performance anti-patterns.
func slowLogProcessor(lines []string) []string {
	var results []string
	for _, line := range lines {
		// ANTI-PATTERN 1: Compiling regex inside a loop is O(N) compilations.
		// The compiled DFA is thrown away after each iteration.
		re := regexp.MustCompile(`ERROR|WARN|FATAL`)
		if re.MatchString(line) {
			// ANTI-PATTERN 2: String concatenation with += is O(N²).
			// Each += allocates a new string and copies the old one.
			result := ""
			result += "[ALERT] "
			result += strings.ToUpper(line)
			results = append(results, result)
		}
	}
	return results
}

// fastLogProcessor fixes both anti-patterns discovered via profiling.
var alertPattern = regexp.MustCompile(`ERROR|WARN|FATAL`) // Compiled ONCE at init

func fastLogProcessor(lines []string) []string {
	results := make([]string, 0, len(lines)/10) // Pre-allocate estimate
	var sb strings.Builder
	for _, line := range lines {
		if alertPattern.MatchString(line) {
			sb.Reset()
			sb.WriteString("[ALERT] ")
			sb.WriteString(strings.ToUpper(line))
			results = append(results, sb.String())
		}
	}
	return results
}

// generateLogs creates a large synthetic log dataset for profiling.
func generateLogs(n int) []string {
	templates := []string{
		"INFO: user %d logged in",
		"ERROR: database connection failed for user %d",
		"WARN: rate limit approaching for user %d",
		"INFO: cache hit for key user:%d",
		"FATAL: out of memory at request %d",
		"DEBUG: query took 42ms for user %d",
	}
	logs := make([]string, n)
	for i := range logs {
		logs[i] = fmt.Sprintf(templates[i%len(templates)], i)
	}
	return logs
}

func main() {
	// =========================================================================
	// Step 1: Start CPU profiling
	// =========================================================================
	cpuFile, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal("could not create CPU profile:", err)
	}
	defer cpuFile.Close()

	// pprof.StartCPUProfile begins sampling. The profiler writes samples to f
	// until StopCPUProfile() is called or the program exits.
	if err := pprof.StartCPUProfile(cpuFile); err != nil {
		log.Fatal("could not start CPU profile:", err)
	}

	// =========================================================================
	// Step 2: Run the workload under profiling
	// =========================================================================
	fmt.Println("Profiling CPU usage...")
	logs := generateLogs(100_000)

	// Run the SLOW version — we expect pprof to show regexp.Compile() as hot
	alerts := slowLogProcessor(logs)
	fmt.Printf("Found %d alerts (slow version)\n", len(alerts))

	// Run the FAST version in the same profile to see the difference
	alerts2 := fastLogProcessor(logs)
	fmt.Printf("Found %d alerts (fast version)\n", len(alerts2))

	// =========================================================================
	// Step 3: Stop CPU profiling
	// =========================================================================
	pprof.StopCPUProfile()
	fmt.Println("CPU profile written to cpu.prof")

	// =========================================================================
	// Step 4: Write Memory profile
	// =========================================================================
	memFile, err := os.Create("mem.prof")
	if err != nil {
		log.Fatal("could not create memory profile:", err)
	}
	defer memFile.Close()

	// Force GC before writing the heap profile — gives cleaner data
	runtime.GC()

	// Lookup "allocs" for allocation profile (vs "heap" for inuse objects)
	if err := pprof.Lookup("allocs").WriteTo(memFile, 0); err != nil {
		log.Fatal("could not write memory profile:", err)
	}
	fmt.Println("Memory profile written to mem.prof")

	// =========================================================================
	// Step 5: What to look for in the profiles
	// =========================================================================
	fmt.Print(`
Next steps:
  go tool pprof -http=:8090 cpu.prof
  → Look for: regexp.Compile(), strings.Builder.copyCheck() in top functions
  → After fix: regexp.Compile() should disappear from the hot path

  go tool pprof -http=:8090 mem.prof
  → Look for: runtime.mallocgc allocations in slowLogProcessor
  → After fix: allocations from fastLogProcessor should be dramatically lower

KEY TAKEAWAY:
  - pprof.StartCPUProfile + StopCPUProfile writes a sampling profile to a file
  - go tool pprof -http=:8090 <file> opens an interactive web UI with flame graph
  - "flat" time = in this function; "cum" time = this function + callees
  - Top anti-patterns found via pprof: regex-in-loop, string +, json on large structs
`)
}
